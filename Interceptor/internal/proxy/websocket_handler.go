package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

var wsConn *websocket.Conn // WebSocket connection to the backend server

// MessageModel represents the structure of the WebSocket message sent to the backend server
type MessageModel struct {
	ApplicationName  string `json:"application_name" binding:"required"`
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

// InitializeWebSocket establishes a WebSocket connection to the backend server
func InitializeWebSocket() error {
	var err error

	// Get the backend host and port from environment variables
	backendHost := os.Getenv("BACKENDHOST")
	if backendHost == "" {
		return fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		return fmt.Errorf("BACKENDPORT environment variable is not set")
	}

	wsURL := fmt.Sprintf("ws://%s:%s/ws", backendHost, backendPort) // Replace with your backend WebSocket server URL
	wsConn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}

	log.Println("WebSocket connection established")
	return nil
}

// SendToBackend sends the log to the backend server through WebSocket
func SendToBackend(message MessageModel) {
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

// CloseWebSocket closes the WebSocket connection
func CloseWebSocket() {
	if wsConn != nil {
		err := wsConn.Close()
		if err != nil {
			log.Println("An error occured while trying to close websocket connection", err)
		}
		log.Println("WebSocket connection closed")
	}
}
