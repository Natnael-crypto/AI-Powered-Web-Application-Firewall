package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRequests(c *gin.Context) {
	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 50
	offset := (page - 1) * pageSize

	var requests []models.Request
	if err := query.Limit(pageSize).Offset(offset).Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func GetRequestByID(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")

	if userRole != "super_admin" {

		userToApplication := models.UserToApplication{
			UserID:          userID,
			ApplicationName: c.Param("application_name"),
		}
		if err := config.DB.Where("user_id = ? AND application_name = ?", userID, c.Param("application_name")).First(&userToApplication).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		requestID := c.Param("request_id")

		var request models.Request
		if err := config.DB.Where("request_id = ?", requestID).First(&request).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}
	}

	requestID := c.Param("request_id")

	var request models.Request
	if err := config.DB.Where("request_id = ?", requestID).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}

func GetAllCountriesStat(c *gin.Context) {

	query := utils.ApplyRequestFilters(c)
	blocked := c.Query("status")

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	query = query.Where("status = ?", blocked)

	type CountryStats struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}

	var stats []CountryStats

	err := query.Select("geo_location as country, count(*) as count").
		Group("geo_location").
		Order("count DESC").
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch country statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocked_countries": stats})
}

func GetRequestsPerMinute(c *gin.Context) {
	blocked := c.Query("blocked")
	var isBlocked bool
	if blocked != "" {
		var err error
		isBlocked, err = strconv.ParseBool(blocked)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid blocked parameter - must be true or false"})
			return
		}
	}

	timeRange := c.DefaultQuery("timerange", "1H")
	rangeValue := timeRange[:len(timeRange)-1]
	rangeUnit := strings.ToUpper(timeRange[len(timeRange)-1:])

	value, err := strconv.Atoi(rangeValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time range value"})
		return
	}

	var duration time.Duration
	switch rangeUnit {
	case "H":
		duration = time.Duration(value) * time.Hour
	case "D":
		duration = time.Duration(value*24) * time.Hour
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time range unit - use H for hours or D for days"})
		return
	}

	intervalStr := c.DefaultQuery("interval", "1")
	intervalMin, err := strconv.Atoi(intervalStr)
	if err != nil || intervalMin <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interval - must be a positive number"})
		return
	}
	intervalDuration := time.Duration(intervalMin) * time.Minute
	intervalSeconds := float64(intervalDuration.Seconds())

	nowUnix := float64(time.Now().Unix())
	startUnix := nowUnix - duration.Seconds()

	query := utils.ApplyRequestFilters(c)
	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	if isBlocked {
		query = query.Where("threat_detected = ?", true)
	}

	type CountResult struct {
		Timestamp float64
		Count     int64
	}

	var results []CountResult

	err = query.
		Where("timestamp >= ?", startUnix).
		Where("timestamp < ?", nowUnix).
		Select("timestamp, COUNT(*) as count").
		Group("timestamp").
		Order("timestamp ASC").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch request counts"})
		return
	}

	buckets := make(map[int64]int64)

	for _, r := range results {
		bucketStart := int64(r.Timestamp/intervalSeconds) * int64(intervalSeconds)
		buckets[bucketStart] += r.Count
	}

	timeSeriesData := make([]map[string]interface{}, 0)
	for ts := int64(startUnix); ts < int64(nowUnix); ts += int64(intervalSeconds) {
		start := ts
		end := ts + int64(intervalSeconds)

		count := buckets[start]

		timeSeriesData = append(timeSeriesData, map[string]interface{}{
			"time_range": fmt.Sprintf("%s - %s",
				time.Unix(start, 0).Format("2006-01-02 15:04:05"),
				time.Unix(end, 0).Format("2006-01-02 15:04:05")),
			"time":  time.Unix(start, 0).Format("2006-01-02 15:04:05"),
			"count": count,
		})
	}

	c.JSON(http.StatusOK, gin.H{"range": timeSeriesData})
}

func GetClientOSStats(c *gin.Context) {

	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	type UserAgentCount struct {
		UserAgent string
		Count     int64
	}

	var uaCounts []UserAgentCount
	err := query.
		Select("user_agent, count(*) as count").
		Where("user_agent IS NOT NULL").
		Group("user_agent").
		Find(&uaCounts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch OS statistics"})
		return
	}

	osStats := map[string]int64{
		"Windows": 0,
		"Linux":   0,
		"macOS":   0,
		"iOS":     0,
		"Android": 0,
		"Other":   0,
	}

	for _, entry := range uaCounts {
		ua := strings.ToLower(entry.UserAgent)
		switch {
		case strings.Contains(ua, "windows"):
			osStats["Windows"] += entry.Count
		case strings.Contains(ua, "linux"):
			osStats["Linux"] += entry.Count
		case strings.Contains(ua, "mac os") || strings.Contains(ua, "macos"):
			osStats["macOS"] += entry.Count
		case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
			osStats["iOS"] += entry.Count
		case strings.Contains(ua, "android"):
			osStats["Android"] += entry.Count
		default:
			osStats["Other"] += entry.Count
		}
	}

	c.JSON(http.StatusOK, gin.H{"os_statistics": osStats})
}

type StatusStats struct {
	ResponseCode int   `json:"response_code"`
	Count        int64 `json:"count"`
}

func GetResponseStatusStats(c *gin.Context) {

	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	query = query.Where("response_code IS NOT NULL")

	var stats []StatusStats
	err := query.
		Select("response_code, COUNT(*) as count").
		Group("response_code").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch response status statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response_status_stats": stats})
}

func GetMostTargetedEndpoints(c *gin.Context) {
	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	query = query.Where("request_url IS NOT NULL AND request_url != ''")

	type EndpointStats struct {
		ApplicationName string `json:"application_name"`
		RequestURL      string `json:"request_url"`
		Count           int64  `json:"count"`
	}

	var stats []EndpointStats

	err := query.
		Select("application_name, request_url, count(*) as count").
		Group("application_name, request_url").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch endpoint statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"most_targeted_endpoints": stats})
}

func GetTopThreatTypes(c *gin.Context) {
	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	query = query.Where("threat_detected = ?", true)

	type AttackStats struct {
		ThreatType string `json:"threat_type"`
		Count      int64  `json:"count"`
	}

	var stats []AttackStats

	err := query.
		Select("threat_type, count(*) as count").
		Group("threat_type").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attack type statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"top_threat_type": stats})
}

func DeleteFilteredRequests(c *gin.Context) {
	query := utils.ApplyRequestFilters(c)
	if query == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to apply filters"})
		return
	}

	result := query.Delete(&models.Request{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Filtered requests deleted successfully",
		"count":   result.RowsAffected,
	})
}

func GetOverallStats(c *gin.Context) {
	query := utils.ApplyRequestFilters(c)

	if query == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to apply filters"})
		return
	}

	// ====== Total Requests ======
	var totalRequests int64
	if err := query.Model(&models.Request{}).Count(&totalRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count total requests"})
		return
	}

	// ====== Blocked Requests & Malicious IPs ======
	blockedQuery := query.Session(&gorm.Session{}) // Clone the query to avoid reuse issues
	blockedQuery = blockedQuery.Where("status = ?", "blocked")

	var blockedRequests int64
	if err := blockedQuery.Count(&blockedRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count blocked requests"})
		return
	}

	var maliciousIPs []string
	if err := blockedQuery.Distinct("client_ip").Pluck("client_ip", &maliciousIPs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count unique malicious IPs"})
		return
	}

	// ====== AI-Based Detections ======
	var aiDetected int64
	if err := blockedQuery.Where("ai_result = ?", true).Count(&aiDetected).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count AI detections"})
		return
	}

	// ====== Rule-Based Detections ======
	var ruleDetected int64
	if err := blockedQuery.Where("rule_detected = ?", true).Count(&ruleDetected).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count rule-based detections"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_requests":        totalRequests,
		"blocked_requests":      blockedRequests,
		"malicious_ips_blocked": len(maliciousIPs),
		"ai_based_detections":   aiDetected,
		"rule_based_detections": ruleDetected,
	})
}
