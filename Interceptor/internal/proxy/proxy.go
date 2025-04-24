package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"interceptor/internal/error_page"
	"interceptor/internal/logger"
	"interceptor/internal/utils"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

var (
	MaintenanceMode bool
	maintenanceLock sync.RWMutex
	ipRateLimiters  = make(map[string]*rate.Limiter)
	limiterLock     sync.Mutex
)

func getTargetRedirectIP(hostname string) (string, bool) {
	appsLock.RLock()
	defer appsLock.RUnlock()
	target, exists := applications[hostname]
	return target, exists
}

func getLimiter(ip string, hostname string) *rate.Limiter {
	limiterLock.Lock()
	defer limiterLock.Unlock()

	if limiter, exists := ipRateLimiters[ip]; exists {
		return limiter
	}

	configLock.RLock()

	limiter := rate.NewLimiter(rate.Limit(application_config[hostname].RateLimit), application_config[hostname].WindowSize)
	configLock.RUnlock()
	ipRateLimiters[ip] = limiter

	fmt.Println(hostname)

	go func() {
		time.Sleep(5 * time.Minute)
		limiterLock.Lock()
		delete(ipRateLimiters, ip)
		limiterLock.Unlock()
	}()

	return limiter
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	maintenanceLock.RLock()
	if MaintenanceMode {
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

	ip := r.RemoteAddr
	ip = strings.Split(ip, ":")[0]
	limiter := getLimiter(ip, r.Host) // TODO check the rate limiting

	if !limiter.Allow() {
		message := utils.MessageModel{
			ApplicationName: r.Host,
			ClientIP:        r.RemoteAddr,
			RequestMethod:   r.Method,
			RequestURL:      r.URL.String(),
			Headers:         fmt.Sprintf("%v", r.Header),
			ResponseCode:    http.StatusTooManyRequests,
			Status:          "blocked",
			ThreatDetected:  true,
			ThreatType:      "Rate Limit Exceeded",
			BotDetected:     false, // TODO Implement bot detection logic
			GeoLocation:     "Unknown",
			RateLimited:     true,
			UserAgent:       r.UserAgent(),
			Body:            "",
		}
		message.Token = WsKey
		utils.SendToBackend(message)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	request_body_size := utils.GetRequestBodySizeMB(r)

	hostname := r.Host

	targetRedirectIP, exists := getTargetRedirectIP(hostname)
	if r.TLS != nil { // TODO targetRedirectIP = "https://" + targetRedirectIP
		targetRedirectIP = "http://" + targetRedirectIP
	} else {
		targetRedirectIP = "http://" + targetRedirectIP
	}
	if !exists {
		http.Error(w, "Unknown host", http.StatusBadGateway)
		return
	}

	wafInstance, exists := wafInstances[hostname]
	if !exists {
		http.Error(w, "WAF instance not found for the application", http.StatusInternalServerError)
		return
	}

	blockedByRule, ruleID, ruleMessage, action, status, body := wafInstance.EvaluateRules(r)

	message := utils.MessageModel{
		ApplicationName: hostname,
		ClientIP:        r.RemoteAddr,
		RequestMethod:   r.Method,
		RequestURL:      r.URL.String(),
		Headers:         fmt.Sprintf("%v", r.Header),
		ResponseCode:    http.StatusOK,
		Status:          fmt.Sprintf("%d", status),
		ThreatDetected:  blockedByRule,
		ThreatType:      ruleMessage,
		BotDetected:     false, // TODO Implement bot detection logic
		GeoLocation:     "Unknown",
		RateLimited:     false,
		UserAgent:       r.UserAgent(),
		Body:            body,
	}

	if request_body_size >= application_config[hostname].MaxPostDataSize {
		if blockedByRule {
			message.Body = utils.HashSHA256(body)
		}
		message.Body = ""
	}

	if blockedByRule {
		error_page.Send403Response(w, ruleID, ruleMessage, action, status)
		logger.LogRequest(r, action, ruleMessage)
		message.ResponseCode = http.StatusForbidden
		message.Status = "blocked"
		message.Token = WsKey

		utils.SendToBackend(message)
		return
	}

	logger.LogRequest(r, "allow", "")
	message.Status = "allowed"

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
	message.Token = WsKey
	utils.SendToBackend(message)
	defer resp.Body.Close()

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
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}
}

func Starter() {
	err := fetchConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = fetchApplications()
	if err != nil {
		log.Fatalf("Failed to fetch applications: %v", err)
	}

	err = utils.InitHttpHandler()
	if err != nil {
		log.Fatalf("Failed to initialize Http Handler: %v", err)
	}

	if remoteLogServer != "" {
		err = logger.InitializeLogger(remoteLogServer)
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}
		defer logger.CloseLogger()
	}

	httpServer := &http.Server{
		Addr:    "0.0.0.0" + proxyPort,
		Handler: http.HandlerFunc(proxyRequest),
	}

	certMap := make(map[string]tls.Certificate)
	CertApp.mu.Lock()
	for _, cert := range CertApp.Certs {
		if cert.CertPath != "" && cert.KeyPath != "" {
			tlsCert, err := tls.LoadX509KeyPair(cert.CertPath, cert.KeyPath)
			if err != nil {
				log.Printf("Failed to load cert for %s: %v", cert.HostName, err)
				continue
			}
			certMap[cert.HostName] = tlsCert
		}
	}
	CertApp.mu.Unlock()

	getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, exists := certMap[chi.ServerName]; exists {
			return &cert, nil
		}
		log.Printf("No certificate found for %s, serving default certificate", chi.ServerName)
		for _, cert := range certMap {
			return &cert, nil
		}
		return nil, fmt.Errorf("no valid certificate found")
	}

	tlsConfig := &tls.Config{
		GetCertificate: getCertificate,
	}

	httpsListener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to create HTTPS listener: %v", err)
	}

	httpsServer := &http.Server{
		Handler: http.HandlerFunc(proxyRequest),
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", proxyPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting HTTPS server on port 443 with SNI support")
		if err := httpsServer.Serve(httpsListener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	if err := httpsServer.Shutdown(ctx); err != nil {
		log.Printf("HTTPS server shutdown error: %v", err)
	}

	log.Println("Servers gracefully stopped")
}
