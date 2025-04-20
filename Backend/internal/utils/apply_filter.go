package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApplyRequestFilters(c *gin.Context) *gorm.DB {
	// Initialize query
	query := config.DB.Model(&models.Request{})

	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	var allowedApps []string

	// 🔹 Role-based Application Filtering
	if userRole != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			return nil
		}

		for _, app := range userApps {
			allowedApps = append(allowedApps, app.ApplicationName)
		}

		if len(allowedApps) == 0 {
			return nil
		}
	}

	// 🔹 Handle application_name filter
	if appFilter := c.QueryArray("application_name"); len(appFilter) > 0 {
		if userRole == "super_admin" {
			query = query.Where("application_name IN ?", appFilter)
		} else {
			var filtered []string
			allowedSet := make(map[string]struct{})
			for _, a := range allowedApps {
				allowedSet[a] = struct{}{}
			}
			for _, a := range appFilter {
				if _, ok := allowedSet[a]; ok {
					filtered = append(filtered, a)
				}
			}

			if len(filtered) == 0 {
				return nil
			}

			query = query.Where("application_name IN ?", filtered)
		}
	} else {
		if userRole != "super_admin" {
			query = query.Where("application_name IN ?", allowedApps)
		}
	}

	// 🔹 Additional Filtering
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

	// 🔹 Date Range Filtering with float-based UNIX timestamps
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

	// 🔹 Specific Date with Time Filtering (expects float UNIX timestamps)
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

	// 🔹 Last X Hours Filtering (based on current UNIX time in float)
	if lastHours := c.Query("last_hours"); lastHours != "" {
		hours, err := strconv.Atoi(lastHours)
		if err != nil {
			return nil
		}
		now := float64(time.Now().Unix())
		past := now - float64(hours*3600)
		query = query.Where("timestamp >= ?", past)
	}

	// 🔹 Full-Text Search
	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	return query
}
