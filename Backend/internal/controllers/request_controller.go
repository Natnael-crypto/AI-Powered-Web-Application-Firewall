package controllers

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddRequest(c *gin.Context) {
	var input struct {
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

	// Parse the input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the application exists
	var app models.Application
	if err := config.DB.Where("application_id = ?", input.ApplicationID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Create the request
	request := models.Request{
		RequestID:        uuid.New().String(),
		ApplicationID:    input.ApplicationID,
		ClientIP:         input.ClientIP,
		RequestMethod:    input.RequestMethod,
		RequestURL:       input.RequestURL,
		Headers:          input.Headers,
		Body:             input.Body,
		Timestamp:        time.Now(),
		ResponseCode:     input.ResponseCode,
		Status:           input.Status,
		MatchedRules:     input.MatchedRules,
		ThreatDetected:   input.ThreatDetected,
		ThreatType:       input.ThreatType,
		ActionTaken:      input.ActionTaken,
		BotDetected:      input.BotDetected,
		GeoLocation:      input.GeoLocation,
		RateLimited:      input.RateLimited,
		UserAgent:        input.UserAgent,
		AIAnalysisResult: input.AIAnalysisResult,
	}

	// Save the request to the database
	if err := config.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "request added successfully", "request": request})
}

// GetRequests retrieves requests based on user role
func GetRequests(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")

	var requests []models.Request

	if userRole == "super_admin" {
		// Super admin: Get all requests
		if err := config.DB.Find(&requests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
			return
		}
	} else {
		// Admin: Get requests for assigned applications
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		// Collect application IDs assigned to the user
		applicationIDs := make([]string, len(userApps))
		for i, app := range userApps {
			applicationIDs[i] = app.ApplicationID
		}

		// Retrieve requests for those applications
		if err := config.DB.Where("application_id IN ?", applicationIDs).Find(&requests).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
