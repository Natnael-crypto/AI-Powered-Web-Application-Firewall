package proxy

// import (
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"interceptor/internal/error_page"
// 	"interceptor/internal/logger"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"os/signal"
// 	"strings"
// 	"sync"
// 	"syscall"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"golang.org/x/time/rate"
// )

// const (
// 	wsReadBufferSize  = 4096
// 	wsWriteBufferSize = 4096
// 	wsMaxMessageSize  = 1024 * 1024
// )

// var (
// 	maintenanceMode bool
// 	maintenanceLock sync.RWMutex
// 	ipRateLimiters  = make(map[string]*rate.Limiter)
// 	limiterLock     sync.Mutex

// 	wsUpgrader = websocket.Upgrader{
// 		ReadBufferSize:  wsReadBufferSize,
// 		WriteBufferSize: wsWriteBufferSize,
// 		CheckOrigin:     func(r *http.Request) bool { return true },
// 	}

// 	wsDialer = &websocket.Dialer{
// 		HandshakeTimeout: 45 * time.Second,
// 		Proxy:            http.ProxyFromEnvironment,
// 	}

// 	wsConnections = struct {
// 		connections map[*websocket.Conn]struct{}
// 		mu          sync.RWMutex
// 	}{
// 		connections: make(map[*websocket.Conn]struct{}),
// 	}
// )

// func isWebSocket(r *http.Request) bool {
// 	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
// 		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request, hostname string) {
// 	maintenanceLock.RLock()
// 	if maintenanceMode {
// 		maintenanceLock.RUnlock()
// 		w.WriteHeader(http.StatusServiceUnavailable)
// 		return
// 	}
// 	maintenanceLock.RUnlock()

// 	ip := strings.Split(r.RemoteAddr, ":")[0]
// 	limiter := getLimiter(ip, hostname)
// 	if !limiter.Allow() {
// 		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
// 		return
// 	}

// 	wafInstance, exists := wafInstances[hostname]
// 	if !exists {
// 		http.Error(w, "WAF configuration not found", http.StatusInternalServerError)
// 		return
// 	}

// 	clientConn, err := wsUpgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("WebSocket upgrade error: %v", err)
// 		return
// 	}
// 	defer clientConn.Close()

// 	wsConnections.mu.Lock()
// 	wsConnections.connections[clientConn] = struct{}{}
// 	wsConnections.mu.Unlock()
// 	defer func() {
// 		wsConnections.mu.Lock()
// 		delete(wsConnections.connections, clientConn)
// 		wsConnections.mu.Unlock()
// 	}()

// 	targetIP, exists := getTargetRedirectIP(hostname)
// 	if !exists {
// 		clientConn.WriteControl(
// 			websocket.CloseMessage,
// 			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Backend not available"),
// 			time.Now().Add(10*time.Second),
// 		)
// 		return
// 	}

// 	backendURL := url.URL{
// 		Scheme: "ws",
// 		Host:   targetIP,
// 		Path:   r.URL.Path,
// 	}
// 	if r.TLS != nil {
// 		backendURL.Scheme = "wss"
// 		wsDialer.TLSClientConfig = &tls.Config{
// 			ServerName:         hostname,
// 			InsecureSkipVerify: true,
// 		}
// 	}

// 	backendConn, _, err := wsDialer.Dial(backendURL.String(), r.Header)
// 	if err != nil {
// 		log.Printf("Backend dial error: %v", err)
// 		clientConn.WriteControl(
// 			websocket.CloseMessage,
// 			websocket.FormatCloseMessage(websocket.CloseServiceRestart, "Cannot connect to backend"),
// 			time.Now().Add(10*time.Second),
// 		)
// 		return
// 	}
// 	defer backendConn.Close()

// 	wsConnections.mu.Lock()
// 	wsConnections.connections[backendConn] = struct{}{}
// 	wsConnections.mu.Unlock()
// 	defer func() {
// 		wsConnections.mu.Lock()
// 		delete(wsConnections.connections, backendConn)
// 		wsConnections.mu.Unlock()
// 	}()

// 	errChan := make(chan error, 2)

// 	go func() {
// 		for {
// 			msgType, message, err := clientConn.ReadMessage()
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}

// 			if blocked, ruleID, _, action, _ := wafInstance.EvaluateWebSocketMessage(message, r); blocked {
// 				logger.LogWebSocketBlock(r, ruleID, action)
// 				clientConn.WriteControl(
// 					websocket.CloseMessage,
// 					websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Blocked by WAF"),
// 					time.Now().Add(10*time.Second),
// 				)
// 				errChan <- fmt.Errorf("message blocked by WAF")
// 				return
// 			}

// 			if err := backendConn.WriteMessage(msgType, message); err != nil {
// 				errChan <- err
// 				return
// 			}
// 		}
// 	}()

// 	go func() {
// 		for {
// 			msgType, message, err := backendConn.ReadMessage()
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}

// 			if err := clientConn.WriteMessage(msgType, message); err != nil {
// 				errChan <- err
// 				return
// 			}
// 		}
// 	}()

// 	select {
// 	case err := <-errChan:
// 		if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 			log.Printf("WebSocket error: %v", err)
// 		}
// 	case <-r.Context().Done():
// 	}
// }

// func getTargetRedirectIP(hostname string) (string, bool) {
// 	appsLock.RLock()
// 	defer appsLock.RUnlock()
// 	target, exists := applications[hostname]
// 	return target, exists
// }

// func getLimiter(ip string, hostname string) *rate.Limiter {
// 	limiterLock.Lock()
// 	defer limiterLock.Unlock()

// 	if limiter, exists := ipRateLimiters[ip]; exists {
// 		return limiter
// 	}

// 	configLock.RLock()
// 	fmt.Printf("Rate Limit: %d, Window Size: %d\n", application_config[hostname].RateLimit, application_config[hostname].WindowSize)
// 	limiter := rate.NewLimiter(rate.Limit(application_config[hostname].RateLimit), application_config[hostname].WindowSize)
// 	configLock.RUnlock()
// 	ipRateLimiters[ip] = limiter

// 	go func() {
// 		time.Sleep(5 * time.Minute)
// 		limiterLock.Lock()
// 		delete(ipRateLimiters, ip)
// 		limiterLock.Unlock()
// 	}()

// 	return limiter
// }

// func StartInterceptor(w http.ResponseWriter, r *http.Request) {
// 	maintenanceLock.Lock()
// 	maintenanceMode = false
// 	maintenanceLock.Unlock()
// 	fmt.Println("Starting interceptor")
// 	w.WriteHeader(http.StatusOK)
// }

// func StopInterceptor(w http.ResponseWriter, r *http.Request) {
// 	maintenanceLock.Lock()
// 	maintenanceMode = true
// 	maintenanceLock.Unlock()
// 	fmt.Println("Stopping interceptor")
// 	w.WriteHeader(http.StatusOK)
// }

// func RestartInterceptor(w http.ResponseWriter, r *http.Request) {
// 	if err := fetchConfig(); err != nil {
// 		log.Fatalf("Error fetching config: %v", err)
// 	}

// 	if err := fetchApplications(); err != nil {
// 		log.Fatalf("Error fetching applications: %v", err)
// 	}
// 	err := InitializeWebSocket()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize WebSocket: %v", err)
// 	}
// 	defer CloseWebSocket()
// 	fmt.Println("Restarting interceptor")
// }

// func proxyRequest(w http.ResponseWriter, r *http.Request) {
// 	if isWebSocket(r) {
// 		handleWebSocket(w, r, r.Host)
// 		return
// 	}

// 	maintenanceLock.RLock()
// 	if maintenanceMode {
// 		maintenanceLock.RUnlock()
// 		w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 		w.WriteHeader(http.StatusServiceUnavailable)
// 		fmt.Fprintf(w, `<!DOCTYPE html>
// 			<html><head><title>Under Maintenance</title></head>
// 			<body><h1>Site Under Maintenance</h1><p>We're currently performing maintenance. Please check back soon.</p></body></html>`)
// 		return
// 	}
// 	maintenanceLock.RUnlock()

// 	if r.URL.Path == "/favicon.ico" {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	ip := r.RemoteAddr
// 	ip = strings.Split(ip, ":")[0]
// 	limiter := getLimiter(ip, r.Host)

// 	if !limiter.Allow() {
// 		message := MessageModel{
// 			ApplicationName:  r.Host,
// 			ClientIP:         r.RemoteAddr,
// 			RequestMethod:    r.Method,
// 			RequestURL:       r.URL.String(),
// 			Headers:          fmt.Sprintf("%v", r.Header),
// 			ResponseCode:     http.StatusTooManyRequests,
// 			Status:           "blocked",
// 			MatchedRules:     "Rate Limit Exceeded",
// 			ThreatDetected:   true,
// 			ThreatType:       "Rate Limit Exceeded",
// 			ActionTaken:      "Rate Limit Exceeded",
// 			BotDetected:      false,
// 			GeoLocation:      "Unknown",
// 			RateLimited:      true,
// 			UserAgent:        r.UserAgent(),
// 			AIAnalysisResult: "",
// 		}
// 		message.Token = WsKey
// 		SendToBackend(message)
// 		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
// 		return
// 	}

// 	hostname := r.Host
// 	targetRedirectIP, exists := getTargetRedirectIP(hostname)
// 	if r.TLS != nil {
// 		targetRedirectIP = "http://" + targetRedirectIP
// 	} else {
// 		targetRedirectIP = "http://" + targetRedirectIP
// 	}
// 	if !exists {
// 		http.Error(w, "Unknown host", http.StatusBadGateway)
// 		return
// 	}

// 	wafInstance, exists := wafInstances[hostname]
// 	if !exists {
// 		http.Error(w, "WAF instance not found for the application", http.StatusInternalServerError)
// 		return
// 	}

// 	blockedByRule, ruleID, ruleMessage, action, status := wafInstance.EvaluateRules(r)

// 	message := MessageModel{
// 		ApplicationName:  hostname,
// 		ClientIP:         r.RemoteAddr,
// 		RequestMethod:    r.Method,
// 		RequestURL:       r.URL.String(),
// 		Headers:          fmt.Sprintf("%v", r.Header),
// 		ResponseCode:     http.StatusOK,
// 		Status:           fmt.Sprintf("%d", status),
// 		MatchedRules:     ruleMessage,
// 		ThreatDetected:   blockedByRule,
// 		ThreatType:       "",
// 		ActionTaken:      action,
// 		BotDetected:      false,
// 		GeoLocation:      "Unknown",
// 		RateLimited:      false,
// 		UserAgent:        r.UserAgent(),
// 		AIAnalysisResult: "",
// 	}

// 	if blockedByRule {
// 		error_page.Send403Response(w, ruleID, ruleMessage, action, status)
// 		logger.LogRequest(r, action, ruleMessage)
// 		message.ResponseCode = http.StatusForbidden
// 		message.Status = "blocked"
// 		message.Token = WsKey

// 		SendToBackend(message)
// 		return
// 	}

// 	logger.LogRequest(r, "allow", "")
// 	message.Status = "allowed"

// 	client := &http.Client{}
// 	targetURL := fmt.Sprintf("%s%s", targetRedirectIP, r.URL.Path)
// 	fmt.Println(targetURL)
// 	if r.URL.RawQuery != "" {
// 		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
// 	}

// 	req, err := http.NewRequest(r.Method, targetURL, r.Body)
// 	if err != nil {
// 		fmt.Printf("Failed to create request: %v", err)
// 		http.Error(w, "Failed to create request", http.StatusInternalServerError)
// 		return
// 	}

// 	for name, values := range r.Header {
// 		for _, value := range values {
// 			req.Header.Add(name, value)
// 		}
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
// 		return
// 	}
// 	message.ResponseCode = resp.StatusCode
// 	message.Token = WsKey
// 	SendToBackend(message)
// 	defer resp.Body.Close()

// 	for name, values := range resp.Header {
// 		for _, value := range values {
// 			w.Header().Add(name, value)
// 		}
// 	}
// 	w.WriteHeader(resp.StatusCode)
// 	_, err = io.Copy(w, resp.Body)
// 	if err != nil {
// 		log.Printf("Failed to copy response body: %v", err)
// 	}
// }

// func Starter() {
// 	err := fetchConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	err = fetchApplications()
// 	if err != nil {
// 		log.Fatalf("Failed to fetch applications: %v", err)
// 	}

// 	err = InitializeWebSocket()
// 	if err != nil {
// 		log.Fatalf("Failed to initialize WebSocket: %v", err)
// 	}
// 	defer CloseWebSocket()

// 	if remoteLogServer != "" {
// 		err = logger.InitializeLogger(remoteLogServer)
// 		if err != nil {
// 			log.Fatalf("Failed to initialize logger: %v", err)
// 		}
// 		defer logger.CloseLogger()
// 	}

// 	httpServer := &http.Server{
// 		Addr:    "0.0.0.0" + proxyPort,
// 		Handler: http.HandlerFunc(proxyRequest),
// 	}

// 	certMap := make(map[string]tls.Certificate)
// 	CertApp.mu.Lock()
// 	for _, cert := range CertApp.Certs {
// 		if cert.CertPath != "" && cert.KeyPath != "" {
// 			tlsCert, err := tls.LoadX509KeyPair(cert.CertPath, cert.KeyPath)
// 			if err != nil {
// 				log.Printf("Failed to load cert for %s: %v", cert.HostName, err)
// 				continue
// 			}
// 			certMap[cert.HostName] = tlsCert
// 		}
// 	}
// 	CertApp.mu.Unlock()

// 	getCertificate := func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 		if cert, exists := certMap[chi.ServerName]; exists {
// 			return &cert, nil
// 		}
// 		log.Printf("No certificate found for %s, serving default certificate", chi.ServerName)
// 		for _, cert := range certMap {
// 			return &cert, nil
// 		}
// 		return nil, fmt.Errorf("no valid certificate found")
// 	}

// 	tlsConfig := &tls.Config{
// 		GetCertificate: getCertificate,
// 	}

// 	httpsListener, err := tls.Listen("tcp", ":443", tlsConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to create HTTPS listener: %v", err)
// 	}

// 	httpsServer := &http.Server{
// 		Handler: http.HandlerFunc(proxyRequest),
// 	}

// 	go func() {
// 		log.Printf("Starting HTTP server on port %s", proxyPort)
// 		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %v", err)
// 		}
// 	}()

// 	go func() {
// 		log.Println("Starting HTTPS server on port 443 with SNI support")
// 		if err := httpsServer.Serve(httpsListener); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTPS server error: %v", err)
// 		}
// 	}()

// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

// 	<-stop
// 	log.Println("Shutting down servers...")

// 	wsConnections.mu.Lock()
// 	for conn := range wsConnections.connections {
// 		conn.WriteControl(
// 			websocket.CloseMessage,
// 			websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"),
// 			time.Now().Add(10*time.Second),
// 		)
// 		conn.Close()
// 		delete(wsConnections.connections, conn)
// 	}
// 	wsConnections.mu.Unlock()

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := httpServer.Shutdown(ctx); err != nil {
// 		log.Printf("HTTP server shutdown error: %v", err)
// 	}

// 	if err := httpsServer.Shutdown(ctx); err != nil {
// 		log.Printf("HTTPS server shutdown error: %v", err)
// 	}

// 	log.Println("Servers gracefully stopped")
// }
