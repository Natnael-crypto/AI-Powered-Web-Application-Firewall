package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProcessBatchRequestService(c *gin.Context) (gin.H, int) {
	var rawRequests []map[string]interface{}
	if err := c.ShouldBindJSON(&rawRequests); err != nil {
		log.Printf("Invalid batch format: %v", err)
		return gin.H{"error": "Invalid batch request"}, http.StatusBadRequest
	}

	appMap := repository.GetApplicationHostMap()
	var requests []models.Request

	for _, requestData := range rawRequests {
		token := safeString(requestData["token"])
		if token != config.WsKey {
			continue
		}

		clientIP := safeString(requestData["client_ip"])
		ipParts := strings.Split(clientIP, ":")
		if len(ipParts) == 0 {
			continue
		}

		headers := utils.ParseHeaders(safeString(requestData["headers"]))
		country := utils.GetCountryName(ipParts[0])

		request := models.Request{
			RequestID:       safeString(requestData["request_id"]),
			ApplicationName: safeString(requestData["application_name"]),
			ApplicationID:   appMap[safeString(requestData["application_name"])],
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
		if err := repository.InsertBatchRequests(requests); err != nil {
			return gin.H{"error": "Database insert failed"}, http.StatusInternalServerError
		}
	}

	return gin.H{"message": "Batch processed successfully", "inserted": len(requests)}, http.StatusOK
}

// Safe casting helpers
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
