package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddRequest(c *gin.Context) {
	var input struct {
		ApplicationName  string `json:"application_name" binding:"required"`
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
		ApplicationName:  input.ApplicationName,
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

// GetRequestStats returns request count and unique IP addresses for specified applications
func GetRequestStats(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	apps := c.QueryArray("apps")
	// Split apps string by comma if provided as single string
	if len(apps) == 1 {
		apps = strings.Split(apps[0], ",")
	}

	fmt.Println(apps)

	// Initialize query builder
	query := config.DB.Model(&models.Request{})

	// Handle application filtering based on role and input
	if userRole == "super_admin" {
		if len(apps) == 1 && apps[0] == "all" {
			// Super admin requesting all apps - no filtering needed
		} else if len(apps) > 0 {
			// Filter by specified applications
			query = query.Where("application_name IN ?", apps)
		}
	} else {
		// For regular admin, get their assigned applications
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		// Extract application names
		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}

		// If specific apps requested, intersect with allowed apps
		if len(apps) > 0 {
			var filteredApps []string
			for _, app := range apps {
				for _, allowedApp := range allowedApps {
					if app == allowedApp {
						filteredApps = append(filteredApps, app)
						break
					}
				}
			}
			query = query.Where("application_name IN ?", filteredApps)
		} else {
			// Use all allowed apps
			query = query.Where("application_name IN ?", allowedApps)
		}
	}

	// Get total request count
	var requestCount int64
	if err := query.Count(&requestCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count requests"})
		return
	}

	// Get unique IP count
	var uniqueIPs []string
	if err := query.Distinct("client_ip").Pluck("client_ip", &uniqueIPs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count unique IPs"})
		return
	}

	// Return stats
	c.JSON(http.StatusOK, gin.H{
		"total_requests": requestCount,
		"unique_ips":     len(uniqueIPs),
	})
}

// GetBlockedStats retrieves statistics about blocked requests
func GetBlockedStats(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	// Parse application names from query params
	apps := c.QueryArray("apps")
	// Split apps string by comma if provided as single string
	if len(apps) == 1 {
		apps = strings.Split(apps[0], ",")
	}

	// Initialize query
	query := config.DB.Model(&models.Request{}).Where("status = ?", "blocked")

	// Handle access control based on role
	if userRole == "super_admin" {
		// Super admin can see all applications
		if len(apps) > 0 {
			query = query.Where("application_name IN ?", apps)
		}
	} else {
		// Regular users can only see their assigned applications
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		// Extract application names
		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}

		// If specific apps requested, intersect with allowed apps
		if len(apps) > 0 {
			var filteredApps []string
			for _, app := range apps {
				for _, allowedApp := range allowedApps {
					if app == allowedApp {
						filteredApps = append(filteredApps, app)
						break
					}
				}
			}
			query = query.Where("application_name IN ?", filteredApps)
		} else {
			// Use all allowed apps
			query = query.Where("application_name IN ?", allowedApps)
		}
	}

	// Get total blocked request count
	var blockedCount int64
	if err := query.Count(&blockedCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count blocked requests"})
		return
	}

	// Get unique blocked IP count
	var uniqueBlockedIPs []string
	if err := query.Distinct("client_ip").Pluck("client_ip", &uniqueBlockedIPs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count unique blocked IPs"})
		return
	}

	// Return stats
	c.JSON(http.StatusOK, gin.H{
		"total_blocked_requests": blockedCount,
		"unique_blocked_ips":     len(uniqueBlockedIPs),
	})
}

// GetTopBlockedCountries returns the top 5 countries with the most blocked requests
func GetTopBlockedCountries(c *gin.Context) {
	// Initialize query
	query := config.DB.Model(&models.Request{}).
		Where("status = ?", "blocked")

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

	type CountryStats struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}

	var stats []CountryStats

	err := query.Select("geo_location as country, count(*) as count").
		Group("geo_location").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch country statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"top_blocked_countries": stats})
}

// GetAllBlockedCountries returns all countries and their blocked request counts
func GetAllBlockedCountries(c *gin.Context) {
	// Initialize query
	query := config.DB.Model(&models.Request{}).
		Where("status = ?", "blocked")

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

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

// GetRequestsPerMinute returns the number of requests per minute for different time intervals
func GetRequestsPerMinute(c *gin.Context) {
	// Get blocked parameter from query
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

	// Get time range from query params (e.g. "2H" for 2 hours, "2D" for 2 days)
	timeRange := c.DefaultQuery("timerange", "1H")

	// Parse the time range value and unit
	rangeValue := timeRange[:len(timeRange)-1]
	rangeUnit := timeRange[len(timeRange)-1:]

	value, err := strconv.Atoi(rangeValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time range value"})
		return
	}

	// Validate time range unit and calculate duration
	var duration time.Duration
	switch strings.ToUpper(rangeUnit) {
	case "H":
		duration = time.Duration(value) * time.Hour
	case "D":
		duration = time.Duration(value*24) * time.Hour
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time range unit - use H for hours or D for days"})
		return
	}

	// Get interval from query params (default to 1 minute if not specified)
	interval := c.DefaultQuery("interval", "1")
	_, err = strconv.Atoi(interval)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interval"})
		return
	}

	// Initialize query
	query := config.DB.Model(&models.Request{})

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

	// Get current time in local timezone
	loc, _ := time.LoadLocation("Local")
	now := time.Now().In(loc)
	startTime := now.Add(-duration)

	// Apply blocked filter if needed
	if isBlocked {
		query = query.Where("threat_detected = ?", true)
	}

	type TimeCount struct {
		TimeInterval time.Time
		Count        int64
	}

	var counts []TimeCount
	err = query.
		Where("timestamp >= ?", startTime).
		Select("date_trunc('minute', timestamp) as time_interval, count(*) as count").
		Group("time_interval").
		Order("time_interval ASC").
		Find(&counts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch request counts"})
		return
	}

	// Create a map for faster lookups
	countMap := make(map[time.Time]int64)
	for _, count := range counts {
		countMap[count.TimeInterval] = count.Count
	}

	// Generate time series data with zero counts for missing minutes
	timeSeriesData := make([]map[string]interface{}, 0)
	currentTime := startTime

	for currentTime.Before(now) || currentTime.Equal(now) {
		count := countMap[currentTime]
		timeSeriesData = append(timeSeriesData, map[string]interface{}{
			"time":  currentTime.Format("2006-01-02 15:04:05"),
			"count": count,
		})
		currentTime = currentTime.Add(time.Minute)
	}

	c.JSON(http.StatusOK, gin.H{"range": timeSeriesData})
}


// GetClientOSStats retrieves statistics about client operating systems from request headers
func GetClientOSStats(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	// Parse application names from query params
	apps := c.QueryArray("apps")
	if len(apps) == 1 {
		apps = strings.Split(apps[0], ",")
	}

	// Initialize query
	query := config.DB.Model(&models.Request{})

	// Handle access control based on role
	if userRole == "super_admin" {
		if len(apps) > 0 {
			query = query.Where("application_name IN ?", apps)
		}
	} else {
		// Regular users can only see their assigned applications
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}

		if len(apps) > 0 {
			var filteredApps []string
			for _, app := range apps {
				for _, allowedApp := range allowedApps {
					if app == allowedApp {
						filteredApps = append(filteredApps, app)
						break
					}
				}
			}
			query = query.Where("application_name IN ?", filteredApps)
		} else {
			query = query.Where("application_name IN ?", allowedApps)
		}
	}

	type OSCount struct {
		OS    string
		Count int64
	}

	var osCounts []OSCount
	err := query.
		Select("headers->>'User-Agent' as os, count(*) as count").
		Where("headers->>'User-Agent' IS NOT NULL").
		Group("headers->>'User-Agent'").
		Find(&osCounts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch OS statistics"})
		return
	}

	// Process User-Agent strings to categorize OS
	osStats := map[string]int64{
		"Windows": 0,
		"Linux":   0,
		"macOS":   0,
		"iOS":     0,
		"Android": 0,
		"Other":   0,
	}

	for _, count := range osCounts {
		userAgent := strings.ToLower(count.OS)
		switch {
		case strings.Contains(userAgent, "windows"):
			osStats["Windows"] += count.Count
		case strings.Contains(userAgent, "linux"):
			osStats["Linux"] += count.Count
		case strings.Contains(userAgent, "mac os") || strings.Contains(userAgent, "macos"):
			osStats["macOS"] += count.Count
		case strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipad"):
			osStats["iOS"] += count.Count
		case strings.Contains(userAgent, "android"):
			osStats["Android"] += count.Count
		default:
			osStats["Other"] += count.Count
		}
	}

	c.JSON(http.StatusOK, gin.H{"os_statistics": osStats})
}


// GetResponseStatusStats returns the top 5 response status codes and their counts
func GetResponseStatusStats(c *gin.Context) {
	// Initialize query
	query := config.DB.Model(&models.Request{})

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

	type StatusStats struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var stats []StatusStats

	err := query.Select("status, count(*) as count").
		Group("status").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch response status statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response_status_stats": stats})
}

// GetMostTargetedEndpoints returns the most frequently targeted endpoints across applications
func GetMostTargetedEndpoints(c *gin.Context) {
	// Initialize query
	query := config.DB.Model(&models.Request{})

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

	type EndpointStats struct {
		ApplicationName string `json:"application_name"`
		Endpoint       string `json:"endpoint"`
		Count         int64  `json:"count"`
	}

	var stats []EndpointStats

	err := query.Select("application_name, endpoint, count(*) as count").
		Group("application_name, endpoint").
		Order("count DESC").
		Limit(10).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch endpoint statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"most_targeted_endpoints": stats})
}


// GetTopAttackTypes returns the top 5 attack types and their counts
func GetTopAttackTypes(c *gin.Context) {
	// Initialize query
	query := config.DB.Model(&models.Request{}).
		Where("threat_detected = ?", true)

	// Get user role and ID from context
	role := c.GetString("role")
	userID := c.GetString("user_id")

	// Handle non-super admin users
	if role != "super_admin" {
		var userApps []models.UserToApplication
		if err := config.DB.Where("user_id = ?", userID).Find(&userApps).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user applications"})
			return
		}

		var allowedApps []string
		for _, ua := range userApps {
			allowedApps = append(allowedApps, ua.ApplicationName)
		}
		query = query.Where("application_name IN ?", allowedApps)
	}

	type AttackStats struct {
		AttackType string `json:"attack_type"`
		Count      int64  `json:"count"`
	}

	var stats []AttackStats

	err := query.Select("attack_type, count(*) as count").
		Group("attack_type").
		Order("count DESC").
		Limit(5).
		Find(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attack type statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"top_attack_types": stats})
}


