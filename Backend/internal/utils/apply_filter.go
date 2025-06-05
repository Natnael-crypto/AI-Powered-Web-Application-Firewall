package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Helper to check for 'not ' prefix
func parseFilterParam(param string) (isNot bool, value string) {
	if strings.HasPrefix(strings.ToLower(param), "not ") {
		return true, strings.TrimSpace(param[4:])
	}
	return false, param
}

func ApplyRequestFilters(c *gin.Context) *gorm.DB {
	query := config.DB.Model(&models.Request{})

	appIDs := GetAssignedApplicationIDs(c)

	if applicationId := c.Query("application_id"); applicationId != "" {
		isNot, value := parseFilterParam(applicationId)
		if slices.Contains(appIDs, value) {
			if isNot {
				query = query.Not("application_id = ?", value)
			} else {
				query = query.Where("application_id = ?", value)
			}
		}
	} else {
		query = query.Where("application_id IN ?", appIDs)
	}

	stringFilters := map[string]string{
		"request_id":       "request_id",
		"client_ip":        "client_ip",
		"application_name": "application_name",
		"request_method":   "request_method",
		"request_url":      "request_url",
		"threat_type":      "threat_type",
		"user_agent":       "user_agent",
		"geo_location":     "geo_location",
		"body":             "body",
		"rule_detected":    "rule_detected",
		"ai_threat_type":   "ai_threat_type",
	}

	for param, column := range stringFilters {
		if val := c.Query(param); val != "" {
			isNot, value := parseFilterParam(val)
			if isNot {
				query = query.Not(column+" ILIKE ?", "%"+value+"%")
			} else {
				query = query.Where(column+" ILIKE ?", "%"+value+"%")
			}
		}
	}

	if val := c.Query("response_code"); val != "" {
		isNot, value := parseFilterParam(val)
		if code, err := strconv.Atoi(value); err == nil {
			if isNot {
				query = query.Not("response_code = ?", code)
			} else {
				query = query.Where("response_code = ?", code)
			}
		}
	}

	boolFilters := map[string]string{
		"ai_result":       "ai_result",
		"threat_detected": "threat_detected",
		"bot_detected":    "bot_detected",
		"rate_limited":    "rate_limited",
	}

	for param, column := range boolFilters {
		if val := c.Query(param); val != "" {
			isNot, value := parseFilterParam(val)
			boolVal := strings.ToLower(value) == "true"
			if isNot {
				query = query.Not(column+" = ?", boolVal)
			} else {
				query = query.Where(column+" = ?", boolVal)
			}
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if startTs, err := strconv.ParseFloat(startDate, 64); err == nil {
			query = query.Where("timestamp >= ?", startTs)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if endTs, err := strconv.ParseFloat(endDate, 64); err == nil {
			query = query.Where("timestamp <= ?", endTs)
		}
	}

	if lastHours := c.Query("last_hours"); lastHours != "" {
		if hours, err := strconv.Atoi(lastHours); err == nil {
			query = query.Where("timestamp >= ?", hours)
		}
	}

	if searchQuery := c.Query("search"); searchQuery != "" {
		query = query.Where("to_tsvector('english', headers || ' ' || body || ' ' || request_url) @@ plainto_tsquery(?)", searchQuery)
	}

	return query
}
