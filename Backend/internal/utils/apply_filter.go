package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApplyRequestFilters(c *gin.Context) *gorm.DB {
	query := config.DB.Model(&models.Request{})

	appIDs := GetAssignedApplicationIDs(c)

	if applicationId := c.Query("application_id"); applicationId != "" {
		if slices.Contains(appIDs, applicationId) {
			query = query.Where("application_id = ?", applicationId)
		}
	} else {
		query = query.Where("application_id IN ?", appIDs)
	}

	if clientIP := c.Query("client_ip"); clientIP != "" {
		query = query.Where("client_ip ILIKE ?", "%"+clientIP+"%")
	}
	if application_name := c.Query("application_name"); application_name != "" {
		query = query.Where("application_name ILIKE ?", "%"+application_name+"%")
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
	if body := c.Query("body"); body != "" {
		query = query.Where("body ILIKE ?", "%"+body+"%")
	}
	if response_code := c.Query("response_code"); response_code != "" {
		query = query.Where("response_code = ?", "%"+response_code+"%")
	}
	if rule_detected := c.Query("rule_detected"); rule_detected != "" {
		query = query.Where("rule_detected ILIKE ?", "%"+rule_detected+"%")
	}
	if ai_threat_type := c.Query("ai_threat_type"); ai_threat_type != "" {
		query = query.Where("ai_threat_type ILIKE ?", "%"+ai_threat_type+"%")
	}
	if c.Query("ai_result") != "" {
		ai_result := c.Query("ai_result") == "true"
		if ai_result {
			query = query.Where("ai_result = ?", ai_result)
		} else {
			query = query.Where("ai_result = ?", false)
		}
	}
	if c.Query("threat_detected") != "" {
		threatDetected := c.Query("threat_detected") == "true"
		if threatDetected {
			query = query.Where("threat_detected = ?", threatDetected)
		} else {
			query = query.Where("threat_detected = ?", false)
		}
	}
	if c.Query("bot_detected") != "" {
		botDetected := c.Query("bot_detected") == "true"
		if botDetected {
			query = query.Where("bot_detected = ?", botDetected)
		} else {
			query = query.Where("bot_detected = ?", false)
		}
	}
	if c.Query("rate_limited") != "" {
		rateLimited := c.Query("rate_limited") == "true"
		if rateLimited {
			query = query.Where("rate_limited = ?", rateLimited)
		} else {
			query = query.Where("rate_limited = ?", false)
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		startTs, err := strconv.ParseFloat(startDate, 64)
		if err == nil {
			query = query.Where("timestamp >= ?", startTs)

		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		endTs, err := strconv.ParseFloat(endDate, 64)
		if err == nil {
			query = query.Where("timestamp <= ?", endTs)

		}
	}

	if lastHours := c.Query("last_hours"); lastHours != "" {
		hours, err := strconv.Atoi(lastHours)
		if err == nil {
			query = query.Where("timestamp >= ?", hours)
		}
	}

	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	return query
}
