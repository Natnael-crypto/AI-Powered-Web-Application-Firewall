package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApplyRequestFilters(c *gin.Context) *gorm.DB {
	query := config.DB.Model(&models.Request{})

	appIDs := GetAssignedApplicationIDs(c)

	if applicationName := c.Query("application_name"); applicationName != "" {
		var application models.Application
		if err := config.DB.Where("application_id In ? And application_name = ?", appIDs, applicationName).First(&application).Error; err != nil {
			log.Fatal(err)
		}
		query = query.Where("application_id = ?", application.ApplicationID)
	} else {
		query = query.Where("application_id IN ?", appIDs)
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

	if startDate := c.Query("start_date"); startDate != "" {
		startTs, err := strconv.ParseFloat(startDate, 64)
		if err != nil {
			return nil
		}
		query = query.Where("timestamp >= ?", startTs)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		endTs, err := strconv.ParseFloat(endDate, 64)
		if err != nil {
			return nil
		}
		query = query.Where("timestamp <= ?", endTs)
	}

	if date := c.Query("date"); date != "" {
		if startTime := c.Query("start_time"); startTime != "" {
			startTs, err := strconv.ParseFloat(startTime, 64)
			if err != nil {
				return nil
			}
			query = query.Where("timestamp >= ?", startTs)
		}
		if endTime := c.Query("end_time"); endTime != "" {
			endTs, err := strconv.ParseFloat(endTime, 64)
			if err != nil {
				return nil
			}
			query = query.Where("timestamp <= ?", endTs)
		}
	}

	if lastHours := c.Query("last_hours"); lastHours != "" {
		hours, err := strconv.Atoi(lastHours)
		if err != nil {
			return nil
		}
		now := float64(time.Now().Unix())
		past := now - float64(hours*3600)
		query = query.Where("timestamp >= ?", past)
	}

	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	return query
}
