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
		ApplicationName    string `json:"application_name" binding:"required"`
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
	if err := config.DB.Where("application_name = ?", input.ApplicationName).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Create the request
	request := models.Request{
		RequestID:        uuid.New().String(),
		ApplicationName:    input.ApplicationName,
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
	query := config.DB.Model(&models.Request{})

	// 🔹 Check if user is an admin (not super_admin), filter by assigned applications
	if userRole != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		// Collect assigned application names
		applicationNames := make([]string, len(userApps))
		for i, app := range userApps {
			applicationNames[i] = app.ApplicationName
		}

		// Restrict query to assigned applications
		if len(applicationNames) > 0 {
			query = query.Where("application_name IN ?", applicationNames)
		} else {
			// If user has no assigned applications, return empty response
			c.JSON(http.StatusOK, gin.H{"requests": []models.Request{}})
			return
		}
	}

	// 🔹 Apply filtering based on query parameters

	if applicationName := c.Query("application_name"); applicationName != "" {
		query = query.Where("application_name ILIKE ?", "%"+applicationName+"%")
	}
	if clientIP := c.Query("client_ip"); clientIP != "" {
		query = query.Where("client_ip ILIKE ?", "%"+clientIP+"%")
	}
	if requestMethod := c.Query("request_method"); requestMethod != "" {
		query = query.Where("request_method ILIKE ?", "%"+requestMethod+"%")
	}
	if requestURL := c.Query("request_url"); requestURL != "" {
		query = query.Where("request_url ILIKE ?", "%"+requestURL+"%")
	}
	if threatType := c.Query("threat_type"); threatType != "" {
		query = query.Where("threat_type ILIKE ?", "%"+threatType+"%")
	}
	if actionTaken := c.Query("action_taken"); actionTaken != "" {
		query = query.Where("action_taken ILIKE ?", "%"+actionTaken+"%")
	}
	if userAgent := c.Query("user_agent"); userAgent != "" {
		query = query.Where("user_agent ILIKE ?", "%"+userAgent+"%")
	}
	if geoLocation := c.Query("geo_location"); geoLocation != "" {
		query = query.Where("geo_location ILIKE ?", "%"+geoLocation+"%")
	}

	// 🔹 Boolean filters
	if c.Query("threat_detected") != "" {
		threatDetected := c.Query("threat_detected") == "true"
		query = query.Where("threat_detected = ?", threatDetected)
	}
	if c.Query("bot_detected") != "" {
		botDetected := c.Query("bot_detected") == "true"
		query = query.Where("bot_detected = ?", botDetected)
	}
	if c.Query("rate_limited") != "" {
		rateLimited := c.Query("rate_limited") == "true"
		query = query.Where("rate_limited = ?", rateLimited)
	}

	// 🔹 Date and Time Filtering
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("timestamp >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("timestamp <= ?", endDate)
	}

	// 🔹 Time Filtering for a Specific Date
	if date := c.Query("date"); date != "" {
		if startTime := c.Query("start_time"); startTime != "" {
			query = query.Where("timestamp >= ?", date+" "+startTime)
		}
		if endTime := c.Query("end_time"); endTime != "" {
			query = query.Where("timestamp <= ?", date+" "+endTime)
		}
	}

	// 🔹 Last X Hours Filtering
	if lastHours := c.Query("last_hours"); lastHours != "" {
		query = query.Where("timestamp >= NOW() - INTERVAL '? HOURS'", lastHours)
	}

	// 🔹 Full-Text Search for Large Fields (headers, body, request_url)
	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	// 🔹 Pagination (Default: 50 results per page)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 50
	offset := (page - 1) * pageSize

	if err := query.Limit(pageSize).Offset(offset).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}
