package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var interceptorClients = make(map[*websocket.Conn]bool) // Store interceptor WebSocket connections

// HandleWebSocket handles WebSocket connections for interceptors
func HandleWebSocket(c *gin.Context) {
	// Upgrade initial HTTP request to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set WebSocket upgrade: %v", err)
		return
	}
	defer conn.Close()

	// Read the first message to verify if the connection is from an interceptor
	// var initialMessage map[string]interface{}
	// if err := conn.ReadJSON(&initialMessage); err != nil {
	// 	log.Printf("Error reading initial message: %v", err)
	// 	return
	// }

	// Check if it's an "IamInterceptor" message
	// clientType, exists := initialMessage["message"]
	// if !exists || clientType != "IamInterceptor" {
	// 	log.Println("Invalid or missing client type message")
	// 	return
	// }

	// Register the interceptor connection
	interceptorClients[conn] = true
	log.Println("Interceptor connected")

	// Listen for messages from the interceptor
	for {
		var message interface{}
		if err := conn.ReadJSON(&message); err != nil {
			log.Printf("WebSocket error: %v", err)
			delete(interceptorClients, conn)
			break
		}

		// Handle the received data from the interceptor
		addRequestFromInterceptor(message)
	}
}

// addRequestFromInterceptor adds the request from an interceptor directly to the database
func addRequestFromInterceptor(message interface{}) {
	// Cast message to the expected structure
	requestData, ok := message.(map[string]interface{})
	if !ok {
		log.Println("Error: Invalid request data received from interceptor")
		return
	}

	matchedRules := ""
	if value, ok := requestData["matched_rules"]; ok {
		if strValue, ok := value.(string); ok {
			matchedRules = strValue
		}
	}

	clientIP, ok := requestData["client_ip"].(string)
	if !ok {
		log.Println("Error: client_ip is not a string")
		return
	}

	ipParts := strings.Split(clientIP, ":")

	country := utils.GetCountryName(ipParts[0])

	// Create a new request from the received data
	request := models.Request{
		RequestID:        uuid.New().String(),
		ApplicationName:  requestData["application_name"].(string),
		ClientIP:         requestData["client_ip"].(string),
		RequestMethod:    requestData["request_method"].(string),
		RequestURL:       requestData["request_url"].(string),
		Headers:          requestData["headers"].(string),
		Body:             requestData["body"].(string),
		Timestamp:        time.Now(),
		ResponseCode:     int(requestData["response_code"].(float64)),
		Status:           requestData["status"].(string),
		MatchedRules:     matchedRules,
		ThreatDetected:   requestData["threat_detected"].(bool),
		ThreatType:       requestData["threat_type"].(string),
		ActionTaken:      requestData["action_taken"].(string),
		BotDetected:      requestData["bot_detected"].(bool),
		GeoLocation:      country,
		RateLimited:      requestData["rate_limited"].(bool),
		UserAgent:        requestData["user_agent"].(string),
		AIAnalysisResult: requestData["ai_analysis_result"].(string),
	}

	// Save the request to the database
	if err := config.DB.Create(&request).Error; err != nil {
		log.Printf("Error saving request from interceptor: %v", err)
		return
	}

	log.Printf("Request saved: %v", request.RequestID)
}
