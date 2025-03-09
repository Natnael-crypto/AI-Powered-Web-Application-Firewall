package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"interceptor/internal/error_page"
	"interceptor/internal/logger"
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
	limiter := getLimiter(ip, r.Host)

	if !limiter.Allow() {
		// Respond with a 429 status code if the rate limit is exceeded
		message := MessageModel{
			ApplicationName:  r.Host,
			ClientIP:         r.RemoteAddr,
			RequestMethod:    r.Method,
			RequestURL:       r.URL.String(),
			Headers:          fmt.Sprintf("%v", r.Header),
			ResponseCode:     http.StatusTooManyRequests,
			Status:           "blocked",
			MatchedRules:     "Rate Limit Exceeded",
			ThreatDetected:   true,
			ThreatType:       "Rate Limit Exceeded",
			ActionTaken:      "Rate Limit Exceeded",
			BotDetected:      false, // TODO Implement bot detection logic
			GeoLocation:      "Unknown",
			RateLimited:      true, // TODO Implement rate limiting logic
			UserAgent:        r.UserAgent(),
			AIAnalysisResult: "", // TODO Replace with actual AI analysis results
		}
		SendToBackend(message)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	hostname := r.Host
	targetRedirectIP, exists := getTargetRedirectIP(hostname)
	if r.TLS != nil {
		targetRedirectIP = "http://" + targetRedirectIP
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
	// Load configuration
	err := fetchConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Fetch applications and configure proxy
	err = fetchApplications()
	if err != nil {
		log.Fatalf("Failed to fetch applications: %v", err)
	}

	err = InitializeWebSocket()
	if err != nil {
		log.Fatalf("Failed to initialize WebSocket: %v", err)
	}
	defer CloseWebSocket()

	err = logger.InitializeLogger(remoteLogServer)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.CloseLogger()

	httpServer := &http.Server{
		Addr:    "0.0.0.0" + proxyPort,
		Handler: http.HandlerFunc(proxyRequest),
	}

	// Load TLS certificates into a map
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

	// Custom certificate selection function
	getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, exists := certMap[chi.ServerName]; exists {
			return &cert, nil
		}
		log.Printf("No certificate found for %s, serving default certificate", chi.ServerName)
		for _, cert := range certMap {
			return &cert, nil // Return the first available cert as a fallback
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

	// Start HTTP and HTTPS servers concurrently
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

	// Graceful shutdown handling
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
