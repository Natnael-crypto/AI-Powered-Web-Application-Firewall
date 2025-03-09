package controllers

import (
	// "backend/internal/controllers"
	"backend/internal/config"
	"backend/internal/models"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateRequestsCSV handles the generation and download of requests data as a CSV file
func GenerateRequestsCSV(c *gin.Context) {
	// Get requests from the database or service
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
		parsedDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
			return
		}
		// Set time to start of day (00:00:00)
		parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
		query = query.Where("timestamp >= ?", parsedDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		parsedDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
			return
		}
		// Set time to end of day (23:59:59)
		parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)
		query = query.Where("timestamp <= ?", parsedDate)
	}

	// 🔹 Time Filtering for a Specific Date
	if date := c.Query("date"); date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}

		loc, _ := time.LoadLocation("Local")

		if startTime := c.Query("start_time"); startTime != "" {
			parsedTime, err := time.Parse("15:04:05", startTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format"})
				return
			}
			// Combine date and time in local timezone
			startDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc)
			query = query.Where("timestamp >= ?", startDateTime)
		}
		if endTime := c.Query("end_time"); endTime != "" {
			parsedTime, err := time.Parse("15:04:05", endTime)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format"})
				return
			}
			// Combine date and time in local timezone
			endDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 999999999, loc)
			query = query.Where("timestamp <= ?", endDateTime)
		}
	}

	// 🔹 Last X Hours Filtering
	if lastHours := c.Query("last_hours"); lastHours != "" {
		hours, err := strconv.Atoi(lastHours)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_hours format"})
			return
		}
		// Use time.Now() in local timezone
		loc, _ := time.LoadLocation("Local")
		now := time.Now().In(loc)
		startTime := now.Add(-time.Duration(hours) * time.Hour)
		query = query.Where("timestamp >= ?", startTime)
	}

	// 🔹 Full-Text Search for Large Fields (headers, body, request_url)
	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	if err := config.DB.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve requests",
		})
		return
	}

	if err := query.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}

	// Set headers for CSV download
	filename := fmt.Sprintf("requests_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	// Create a CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write CSV header
	headers := []string{
		"ID", "Application Name", "Client IP", "Request Method", "Request URL", "Headers", "Body", "Timestamp", "ResponseCode", "Status", "MatchedRules", "ThreatType", "ActionTaken", "BotDetected", "GeoLocation", "RateLimited", "UserAgent", "AIAnalysisResult",
		// Add more fields as needed based on your request structure
	}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to write CSV headers",
		})
		return
	}

	// Write request data rows
	for _, req := range requests {
		row := []string{
			req.RequestID,
			req.ApplicationName,
			req.ClientIP,
			req.RequestMethod,
			req.RequestURL,
			req.Headers,
			req.Body,
			req.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%d", req.ResponseCode),
			req.Status,
			req.MatchedRules,
			req.ThreatType,
			fmt.Sprintf("%t", req.BotDetected),
			req.GeoLocation,
			fmt.Sprintf("%t", req.RateLimited),
			req.UserAgent,
			req.AIAnalysisResult,
			// Add more fields as needed
		}
		if err := writer.Write(row); err != nil {
			// Log the error but continue processing
			fmt.Printf("Error writing request %s to CSV: %v\n", req.RequestID, err)
		}
	}

	// The writer will write to the response when Flush() is called
}
