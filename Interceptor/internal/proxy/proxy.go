package proxy

import (
	"fmt"
	"interceptor/internal/error_page"
	"interceptor/internal/logger"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	maintenanceMode bool
	maintenanceLock sync.RWMutex
	ipRateLimiters  = make(map[string]*rate.Limiter)
	limiterLock     sync.Mutex
)

// getTargetRedirectIP gets the backend target URL based on request hostname
func getTargetRedirectIP(hostname string) (string, bool) {
	appsLock.RLock()
	defer appsLock.RUnlock()
	target, exists := applications[hostname]
	return target, exists
}

// getLimiter retrieves or creates a new rate limiter for each client IP
func getLimiter(ip string) *rate.Limiter {
	limiterLock.Lock()
	defer limiterLock.Unlock()

	if limiter, exists := ipRateLimiters[ip]; exists {
		return limiter
	}

	// Create a new limiter: 5 requests per second with a burst of 10
	// limiter := rate.NewLimiter(1, 10)
	configLock.RLock()
	limiter := rate.NewLimiter(rate.Limit(rateLimit), windowSize)
	configLock.RUnlock()
	ipRateLimiters[ip] = limiter

	// Optional: clean up old entries periodically (could use a background go routine)
	go func() {
		time.Sleep(5 * time.Minute)
		limiterLock.Lock()
		delete(ipRateLimiters, ip)
		limiterLock.Unlock()
	}()

	return limiter
}

func StartInterceptor(w http.ResponseWriter, r *http.Request) {
	maintenanceLock.Lock()
	maintenanceMode = false
	maintenanceLock.Unlock()
	fmt.Println("Starting interceptor")
	w.WriteHeader(http.StatusOK)
}

func StopInterceptor(w http.ResponseWriter, r *http.Request) {
	maintenanceLock.Lock()
	maintenanceMode = true
	maintenanceLock.Unlock()
	fmt.Println("Stopping interceptor")
	w.WriteHeader(http.StatusOK)
}

func RestartInterceptor(w http.ResponseWriter, r *http.Request) {
	if err := fetchConfig(); err != nil {
		log.Fatalf("Error fetching config: %v", err)
	}

	if err := fetchApplications(); err != nil {
		log.Fatalf("Error fetching applications: %v", err)
	}
	err := InitializeWebSocket()
	if err != nil {
		log.Fatalf("Failed to initialize WebSocket: %v", err)
	}
	defer CloseWebSocket()
	fmt.Println("Restarting interceptor")
}

// proxyRequest handles incoming requests and forwards them to the correct backend server
func proxyRequest(w http.ResponseWriter, r *http.Request) {
	maintenanceLock.RLock()
	if maintenanceMode {
		maintenanceLock.RUnlock()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, `<!DOCTYPE html>
			<html><head><title>Under Maintenance</title></head>
			<body><h1>Site Under Maintenance</h1><p>We're currently performing maintenance. Please check back soon.</p></body></html>`)
		return
	}
	maintenanceLock.RUnlock()

	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Apply rate limiting here
	ip := r.RemoteAddr
	ip = strings.Split(ip, ":")[0]
	limiter := getLimiter(ip)

	if !limiter.Allow() {
		// Respond with a 429 status code if the rate limit is exceeded
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	hostname := r.Host
	targetRedirectIP, exists := getTargetRedirectIP(hostname)
	if r.TLS != nil {
		targetRedirectIP = "https://" + targetRedirectIP
	} else {
		targetRedirectIP = "http://" + targetRedirectIP
	}
	if !exists {
		http.Error(w, "Unknown host", http.StatusBadGateway)
		return
	}

	// Retrieve the WAF instance for the application based on the hostname
	wafInstance, exists := wafInstances[hostname]
	if !exists {
		http.Error(w, "WAF instance not found for the application", http.StatusInternalServerError)
		return
	}

	// Evaluate the request using the appropriate WAF instance
	blockedByRule, ruleID, ruleMessage, action, status := wafInstance.EvaluateRules(r)
	message := MessageModel{
		ApplicationName:  hostname,
		ClientIP:         r.RemoteAddr,
		RequestMethod:    r.Method,
		RequestURL:       r.URL.String(),
		Headers:          fmt.Sprintf("%v", r.Header),
		ResponseCode:     http.StatusOK,
		Status:           fmt.Sprintf("%d", status),
		MatchedRules:     ruleMessage,
		ThreatDetected:   blockedByRule,
		ThreatType:       "", // TODO Replace with actual threat detection logic
		ActionTaken:      action,
		BotDetected:      false, // TODO Implement bot detection logic
		GeoLocation:      "Unknown",
		RateLimited:      false, // TODO Implement rate limiting logic
		UserAgent:        r.UserAgent(),
		AIAnalysisResult: "", // TODO Replace with actual AI analysis results
	}

	if blockedByRule {
		error_page.Send403Response(w, ruleID, ruleMessage, action, status)
		logger.LogRequest(r, action, ruleMessage)
		message.ResponseCode = http.StatusForbidden
		message.Status = "blocked"
		SendToBackend(message) // Send log to backend
		return
	}

	logger.LogRequest(r, "allow", "")
	message.Status = "allowed"
	// Send log to backend

	client := &http.Client{}
	targetURL := fmt.Sprintf("%s%s", targetRedirectIP, r.URL.Path)
	fmt.Println(targetURL)
	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		fmt.Printf("Failed to create request: %v", err)
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
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		return
	}
	message.ResponseCode = resp.StatusCode
	SendToBackend(message)
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}
}

func Starter() {
	if err := fetchConfig(); err != nil {
		log.Fatalf("Error fetching config: %v", err)
	}

	if err := fetchApplications(); err != nil {
		log.Fatalf("Error fetching applications: %v", err)
	}

	// Initialize WebSocket connection
	err := InitializeWebSocket()
	if err != nil {
		log.Fatalf("Failed to initialize WebSocket: %v", err)
	}
	defer CloseWebSocket()

	http.HandleFunc("/", proxyRequest)
	http.HandleFunc("/interceptor/startinterceptor", StartInterceptor)
	http.HandleFunc("/interceptor/stopinterceptor", StopInterceptor)
	http.HandleFunc("/interceptor/restartinterceptor", RestartInterceptor)

	err = logger.InitializeLogger(remoteLogServer)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.CloseLogger()

	fmt.Printf("Starting server on port %s\n", proxyPort)
	if err := http.ListenAndServe("0.0.0.0"+proxyPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
