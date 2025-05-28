package proxy

import (
	"fmt"
	"interceptor/internal/error_page"
	"interceptor/internal/fusionService"
	"interceptor/internal/logger"
	"interceptor/internal/ml"
	"interceptor/internal/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

var (
	MaintenanceMode bool
	maintenanceLock sync.RWMutex
	ipRateLimiters  = make(map[string]*rate.Limiter)
	blockedIPs      = make(map[string]time.Time)
	limiterLock     sync.Mutex
)

func getTargetRedirectIP(hostname string) (string, bool) {
	appsLock.RLock()
	defer appsLock.RUnlock()
	target, exists := applications[hostname]
	return target, exists
}

func getLimiter(ip, hostname string) *rate.Limiter {
	limiterLock.Lock()
	defer limiterLock.Unlock()

	if limiter, exists := ipRateLimiters[ip]; exists {
		return limiter
	}

	configLock.RLock()
	limiter := rate.NewLimiter(rate.Limit(application_config[hostname].RateLimit), application_config[hostname].WindowSize)
	configLock.RUnlock()
	ipRateLimiters[ip] = limiter

	return limiter
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	// Maintenance mode check
	maintenanceLock.RLock()
	if MaintenanceMode {
		maintenanceLock.RUnlock()
		error_page.SendMaintenanceResponse(w)
		return
	}
	maintenanceLock.RUnlock()

	ip := strings.Split(r.RemoteAddr, ":")[0]
	hostname := r.Host

	// Rate limiter check
	limiterLock.Lock()
	if unblockTime, blocked := blockedIPs[ip]; blocked {
		if time.Now().Before(unblockTime) {
			limiterLock.Unlock()
			http.Error(w, "Too Many Requests (Blocked)", http.StatusTooManyRequests)
			logger.LogRequest(r, "blocked", hostname, ip, 0.6)

			message := utils.MessageModel{
				RequestID:       uuid.New().String(),
				ApplicationName: hostname,
				ClientIP:        r.RemoteAddr,
				RequestMethod:   r.Method,
				RequestURL:      utils.RecursiveDecode(r.URL.String(), 3),
				Headers:         fmt.Sprintf("%v", r.Header),
				ResponseCode:    http.StatusTooManyRequests,
				Status:          "blocked",
				ThreatDetected:  true,
				ThreatType:      "Rate Limit Exceeded",
				BotDetected:     false,
				GeoLocation:     "Unknown",
				RateLimited:     true,
				UserAgent:       r.UserAgent(),
				Body:            "",
				Token:           WsKey,
				AIResult:        false,
				RuleDetected:    false,
				AIThreatType:    "",
			}
			utils.SendToBackend(message)

			return
		}
		delete(blockedIPs, ip)
	}
	limiterLock.Unlock()

	limiter := getLimiter(ip, hostname)
	if !limiter.Allow() {
		limiterLock.Lock()
		blockedIPs[ip] = time.Now().Add(time.Duration(application_config[hostname].BlockTime) * time.Minute)
		limiterLock.Unlock()

		message := utils.MessageModel{
			RequestID:       uuid.New().String(),
			ApplicationName: hostname,
			ClientIP:        r.RemoteAddr,
			RequestMethod:   r.Method,
			RequestURL:      utils.RecursiveDecode(r.URL.String(), 3),
			Headers:         fmt.Sprintf("%v", r.Header),
			ResponseCode:    http.StatusTooManyRequests,
			Status:          "blocked",
			ThreatDetected:  true,
			ThreatType:      "Rate Limit Exceeded",
			BotDetected:     false,
			GeoLocation:     "Unknown",
			RateLimited:     true,
			UserAgent:       r.UserAgent(),
			Body:            "",
			Token:           WsKey,
			AIResult:        false,
			RuleDetected:    false,
			AIThreatType:    "",
		}
		utils.SendToBackend(message)
		logger.LogRequest(r, "blocked", hostname, ip, 0.6)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// Target redirect IP resolution
	targetRedirectIP, exists := getTargetRedirectIP(hostname)
	if !exists {
		http.Error(w, "Unknown host", http.StatusBadGateway)
		logger.LogRequest(r, "blocked", hostname, ip, 0.4)
		return
	}

	if application_config[hostname].Tls {
		targetRedirectIP = "https://" + targetRedirectIP
	} else {
		targetRedirectIP = "http://" + targetRedirectIP
	}

	wafInstance, exists := wafInstances[hostname]
	if !exists {
		http.Error(w, "WAF instance not found for the application", http.StatusInternalServerError)
		logger.LogRequest(r, "blocked", hostname, ip, 0.4)
		return
	}

	// WAF rule evaluation
	blockedByRule, ruleID, ruleMessage, _, status, body := wafInstance.EvaluateRules(r)
	headers := utils.ParseHeaders(fmt.Sprintf("%v", r.Header))
	requestBodySize := utils.GetRequestBodySizeMB(r)

	// Prepare message
	message := utils.MessageModel{
		RequestID:       uuid.New().String(),
		ApplicationName: hostname,
		ClientIP:        r.RemoteAddr,
		RequestMethod:   r.Method,
		RequestURL:      utils.RecursiveDecode(r.URL.String(), 3),
		Headers:         fmt.Sprintf("%v", r.Header),
		ResponseCode:    http.StatusOK,
		Status:          fmt.Sprintf("%d", status),
		ThreatDetected:  blockedByRule,
		ThreatType:      ruleMessage,
		BotDetected:     false,
		GeoLocation:     "Unknown",
		RateLimited:     false,
		UserAgent:       r.UserAgent(),
		Body:            body,
		AIResult:        false,
		RuleDetected:    false,
		AIThreatType:    "",
	}

	// Special WAF Rule ID - Skip ML & Fusion
	if blockedByRule && ruleID >= 1000000000000000000 {
		if requestBodySize >= application_config[hostname].MaxPostDataSize {
			message.Body = ""
		} else {
			message.Body = utils.HashSHA256(body)
		}

		error_page.Send403Response(w, message.RequestID)
		message.ResponseCode = http.StatusForbidden
		message.Status = "blocked"
		message.Token = WsKey
		message.RuleDetected = true
		logger.LogRequest(r, "blocked", hostname, ip, 0.6)
		utils.SendToBackend(message)
		return
	}

	// ML and Fusion Evaluation
	requestData := ml.RequestData{
		RequestID: message.RequestID,
		Url:       r.URL.String(),
		Headers:   headers,
		Body:      body,
	}

	blockedByMl, Normal, Anomaly, err := ml.EvaluateML(requestData)
	// fmt.Print(blockedByMl, Normal, Anomaly)
	if err != nil {
		http.Error(w, "Error evaluating ML model", http.StatusInternalServerError)
		logger.LogRequest(r, "blocked", hostname, ip, 0.4)
		return
	}

	message.AIResult = blockedByMl
	message.RuleDetected = blockedByRule

	result := fusionService.FusionAlgorithm(blockedByRule, Normal, Anomaly)

	if requestBodySize >= application_config[hostname].MaxPostDataSize {
		if result {
			message.Body = utils.HashSHA256(body)
		}
		message.Body = ""
	}

	// err = utils.SaveEvaluationResult(message.RuleDetected, Normal, Anomaly, "results_n.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if result {
		error_page.Send403Response(w, message.RequestID)
		message.ResponseCode = http.StatusForbidden
		message.Status = "blocked"
		message.Token = WsKey
		utils.SendToBackend(message)
		logger.LogRequest(r, "blocked", hostname, ip, 0.6)
		return
	}

	// Proxy forwarding
	message.Status = "allowed"
	client := &http.Client{}
	targetURL := fmt.Sprintf("%s%s", targetRedirectIP, r.URL.Path)
	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		logger.LogRequest(r, "allowed", hostname, ip, 0.4)
		return
	}
	defer resp.Body.Close()

	message.ResponseCode = resp.StatusCode
	message.Token = WsKey
	utils.SendToBackend(message)

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	if headers, ok := application_security_headers[hostname]; ok {
		for _, h := range headers {
			w.Header().Set(h.HeaderName, h.HeaderValue)
		}
	}

	w.WriteHeader(resp.StatusCode)
	logger.LogRequest(r, "allowed", hostname, ip, 0.4)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}
}
