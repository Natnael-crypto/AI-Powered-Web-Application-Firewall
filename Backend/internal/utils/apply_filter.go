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

	// 🔹 Date Range Filtering
	if startDate := c.Query("start_date"); startDate != "" {
		parsedDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil
		}
		parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
		query = query.Where("timestamp >= ?", parsedDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		parsedDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil
		}
		parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)
		query = query.Where("timestamp <= ?", parsedDate)
	}

	// 🔹 Specific Date with Time Filtering
	if date := c.Query("date"); date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil
		}
		loc, _ := time.LoadLocation("Local")

		if startTime := c.Query("start_time"); startTime != "" {
			parsedTime, err := time.Parse("15:04:05", startTime)
			if err != nil {
				return nil
			}
			startDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, loc)
			query = query.Where("timestamp >= ?", startDateTime)
		}
		if endTime := c.Query("end_time"); endTime != "" {
			parsedTime, err := time.Parse("15:04:05", endTime)
			if err != nil {
				return nil
			}
			endDateTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 999999999, loc)
			query = query.Where("timestamp <= ?", endDateTime)
		}
	}

	// 🔹 Last X Hours Filtering
	if lastHours := c.Query("last_hours"); lastHours != "" {
		hours, err := strconv.Atoi(lastHours)
		if err != nil {
			return nil
		}
		loc, _ := time.LoadLocation("Local")
		now := time.Now().In(loc)
		startTime := now.Add(-time.Duration(hours) * time.Hour)
		query = query.Where("timestamp >= ?", startTime)
	}

	// 🔹 Full-Text Search
	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	return query
}
