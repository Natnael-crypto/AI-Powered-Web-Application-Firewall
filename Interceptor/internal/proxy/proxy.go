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
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

const remoteLogServer = "192.168.1.2:514" // TODO Get the IP address and port from the server

const targetRedirectIP = "http://127.0.0.1:5500" // TODO Get the IP address and port from the server

const proxyPort = ":8000" //TODO Get the port from the server

var requestQueue = make(chan *http.Request, 100)
var wsConn *websocket.Conn // WebSocket connection to the backend server

// MessageModel represents the structure of the WebSocket message sent to the backend server
type MessageModel struct {
	ApplicationID    string `json:"application_id" binding:"required"`
	ClientIP         string `json:"client_ip" binding:"required"`
	RequestMethod    string `json:"request_method" binding:"required"`
	RequestURL       string `json:"request_url" binding:"required"`
	Headers          string `json:"headers"`
	Body             string `json:"body"`
	ResponseCode     int    `json:"response_code" binding:"required"`
	Status           string `json:"status" binding:"required"`
	MatchedRules     string `json:"matched_rules"`
	ThreatDetected   bool   `json:"threat_detected"`
	ThreatType       string `json:"threat_type"`
	ActionTaken      string `json:"action_taken"`
	BotDetected      bool   `json:"bot_detected"`
	GeoLocation      string `json:"geo_location"`
	RateLimited      bool   `json:"rate_limited"`
	UserAgent        string `json:"user_agent"`
	AIAnalysisResult string `json:"ai_analysis_result"`
}

// sendToBackend sends the log to the backend server through WebSocket
func sendToBackend(message MessageModel) {
	if wsConn == nil {
		log.Println("WebSocket connection is not established")
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	err = wsConn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("Failed to send message through WebSocket: %v", err)
	}
}

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusOK)
		return
	}

	blocked_by_Rule, RuleID, RuleMessage, Action, Status := waf.EvaluateRules(r)

	message := MessageModel{
		ApplicationID:    "b028fd26-fd2c-4486-8a1f-d1b510b652f0",
		ClientIP:         r.RemoteAddr,
		RequestMethod:    r.Method,
		RequestURL:       r.URL.String(),
		Headers:          fmt.Sprintf("%v", r.Header),
		ResponseCode:     http.StatusOK,
		Status:           fmt.Sprintf("%d", Status),
		MatchedRules:     RuleMessage,
		ThreatDetected:   blocked_by_Rule,
		ThreatType:       "SQL Injection", // TODO Replace this with the actual threat type logic
		ActionTaken:      Action,
		BotDetected:      false,     // TODO Add logic to detect bot traffic if required
		GeoLocation:      "Unknown", //TODO Add logic to resolve GeoLocation
		RateLimited:      false,     //TODO Add logic for rate limiting if required
		UserAgent:        r.UserAgent(),
		AIAnalysisResult: "No malicious activity detected", // TODO Replace with actual AI analysis results
	}

	if blocked_by_Rule {
		error_page.Send403Response(w, RuleID, RuleMessage, Action, Status)
		logger.LogRequest(r, Action, RuleMessage)
		message.ResponseCode = http.StatusForbidden
		message.Status = "blocked"
		sendToBackend(message) // Send log to backend
		return
	}

	logger.LogRequest(r, "allow", "")
	message.Status = "allowed"
	sendToBackend(message) // Send log to backend

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

func worker() {
	for r := range requestQueue {
		fmt.Printf("Worker processing request: %s %s\n", r.Method, r.URL)
	}
}

func initializeWebSocket() error {
	var err error
	wsURL := "ws://localhost:8080/ws" // Replace with your backend WebSocket server URL
	wsConn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket server: %v", err)
	}

	log.Println("WebSocket connection established")
	return nil
}

func Starter() {
	if err := waf.InitializeRuleEngine(); err != nil {
		log.Fatalf("Failed to initialize WAF: %v", err)
	}

	workerPoolSize := runtime.NumCPU()
	fmt.Printf("Detected %d logical CPUs. Setting worker pool size to %d.\n", workerPoolSize, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		go worker()
	}

	// Initialize WebSocket connection
	err := initializeWebSocket()
	if err != nil {
		log.Fatalf("Failed to initialize WebSocket: %v", err)
	}
	defer func() {
		if wsConn != nil {
			err := wsConn.Close()
			if err != nil {
				log.Println("An error occurred trying to close websocket connection", err)
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case requestQueue <- r:
		default:
			http.Error(w, "Server too busy. Try again later.", http.StatusServiceUnavailable)
		}
		proxyRequest(w, r)
	})

	err = logger.InitializeLogger(remoteLogServer)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.CloseLogger()

	fmt.Printf("Starting server on port %s\n", proxyPort)
	server := &http.Server{
		Addr:              ":1234",
		ReadHeaderTimeout: 3 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
