package proxy

import (
	"encoding/json"
	"fmt"
	"interceptor/internal/error_page"
	"interceptor/internal/logger"
	"interceptor/internal/waf"
	"io"
	"log"
	"net/http"
	"sync"
)

type Config struct {
	ID              string `json:"ID"`
	ListeningPort   string `json:"ListeningPort"`
	RemoteLogServer string `json:"RemoteLogServer"`
}

type Application struct {
	ApplicationID   string `json:"application_id"`
	ApplicationName string `json:"application_name"`
	Hostname        string `json:"hostname"`
	IPAddress       string `json:"ip_address"`
	Port            string `json:"port"`
	Status          bool   `json:"status"`
}

var (
	remoteLogServer string
	proxyPort       string
	applications    map[string]string   // Maps hostname -> "IP:Port"
	wafInstances    map[string]*waf.WAF // Maps hostname -> WAF instance
	configLock      sync.RWMutex
	appsLock        sync.RWMutex
)

// fetchConfig retrieves the configuration from the remote API
func fetchConfig() error {
	resp, err := http.Get("http://localhost:8080/config")
	if err != nil {
		return fmt.Errorf("failed to fetch config: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Config Config `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode config: %v", err)
	}

	configLock.Lock()
	remoteLogServer = result.Config.RemoteLogServer
	proxyPort = ":" + result.Config.ListeningPort
	configLock.Unlock()

	fmt.Printf("Loaded config: LogServer=%s, ProxyPort=%s\n", remoteLogServer, proxyPort)
	return nil
}

// fetchApplications retrieves the list of applications and updates the hostname mapping
func fetchApplications() error {
	resp, err := http.Get("http://localhost:8080/application/")
	if err != nil {
		return fmt.Errorf("failed to fetch applications: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Applications []Application `json:"applications"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode applications: %v", err)
	}

	appsLock.Lock()
	applications = make(map[string]string)
	wafInstances = make(map[string]*waf.WAF)
	for _, app := range result.Applications {
		if app.Status {
			applications[app.ApplicationName] = app.Hostname

			// Initialize WAF for each application using its application ID
			rulesResponse, err := FetchRules(app.ApplicationID)
			if err != nil {
				log.Fatalf("Error fetching rules: %v", err)
			}

			// Write the rule strings to a file and get the filename
			fileName, err := WriteRuleToFile(app.ApplicationID, rulesResponse.Rules)
			if err != nil {
				log.Fatalf("Error writing rules to file: %v", err)
			}

			wafInstance, err := waf.InitializeRuleEngine(fileName)
			if err != nil {
				log.Printf("Error initializing WAF for application %s: %v", app.ApplicationName, err)
				continue
			}

			// Store the WAF instance by hostname
			wafInstances[app.ApplicationName] = wafInstance
		}
	}
	appsLock.Unlock()

	fmt.Println("Loaded applications:", applications)
	return nil
}

// getTargetRedirectIP gets the backend target URL based on request hostname
func getTargetRedirectIP(hostname string) (string, bool) {
	appsLock.RLock()
	defer appsLock.RUnlock()
	target, exists := applications[hostname]
	return target, exists
}

// proxyRequest handles incoming requests and forwards them to the correct backend server
func proxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusOK)
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
		BotDetected:      false,     // TODO Implement bot detection logic
		GeoLocation:      "Unknown",
		RateLimited:      false,     // TODO Implement rate limiting logic
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
	SendToBackend(message) // Send log to backend

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

	err = logger.InitializeLogger(remoteLogServer)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.CloseLogger()

	fmt.Printf("Starting server on port %s\n", proxyPort)
	if err := http.ListenAndServe(proxyPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
