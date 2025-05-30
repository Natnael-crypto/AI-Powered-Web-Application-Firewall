package controllers

import (
	"log"
	"net/http"
	"strings"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func safeString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func safeFloat64(val interface{}) float64 {
	if f, ok := val.(float64); ok {
		return f
	}
	return 0
}

func safeBool(val interface{}) bool {
	if b, ok := val.(bool); ok {
		return b
	}
	return false
}

func HandleBatchRequests(c *gin.Context) {
	var rawRequests []map[string]interface{}

	if err := c.ShouldBindJSON(&rawRequests); err != nil {
		log.Printf("Invalid batch format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch request"})
		return
	}

	var requests []models.Request

	var list_of_apps = make(map[string]string)

	var applications []models.Application

	if err := config.DB.Find(&applications).Error; err != nil {
		log.Println("Application not found")
	}

	for _, app := range applications {
		list_of_apps[app.HostName] = app.ApplicationID
	}

	for _, requestData := range rawRequests {
		token, ok := requestData["token"].(string)
		if !ok || token != config.WsKey {
			log.Println("Unauthorized or missing token in batch item")
			continue
		}

		clientIP, ok := requestData["client_ip"].(string)
		if !ok {
			log.Println("Invalid or missing client_ip")
			continue
		}

		ipParts := strings.Split(clientIP, ":")
		country := utils.GetCountryName(ipParts[0])
		headers := utils.ParseHeaders(safeString(requestData["headers"]))

		request := models.Request{
			RequestID:       uuid.New().String(),
			ApplicationName: safeString(requestData["application_name"]),
			ApplicationID:   list_of_apps[safeString(requestData["application_name"])],
			ClientIP:        ipParts[0],
			RequestMethod:   safeString(requestData["request_method"]),
			RequestURL:      safeString(requestData["request_url"]),
			Headers:         headers,
			Body:            safeString(requestData["body"]),
			Timestamp:       safeFloat64(requestData["timestamp"]),
			ResponseCode:    int(safeFloat64(requestData["response_code"])),
			Status:          safeString(requestData["status"]),
			ThreatDetected:  safeBool(requestData["threat_detected"]),
			ThreatType:      safeString(requestData["threat_type"]),
			BotDetected:     safeBool(requestData["bot_detected"]),
			GeoLocation:     country,
			RateLimited:     safeBool(requestData["rate_limited"]),
			UserAgent:       safeString(requestData["user_agent"]),
			AIResult:        safeBool(requestData["ai_result"]),
			AIThreatType:    safeString(requestData["ai_threat_type"]),
			RuleDetected:    safeBool(requestData["rule_detected"]),
		}

		requests = append(requests, request)
	}

	if len(requests) > 0 {
		if err := config.DB.Create(&requests).Error; err != nil {
			log.Printf("Bulk DB insert error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert failed"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch processed successfully", "inserted": len(requests)})
}
