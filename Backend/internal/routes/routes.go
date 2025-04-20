package routes

import (
	"backend/internal/controllers"
	middleware "backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes initializes all the routes for the application
func InitializeRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired)

	// User Management
	authorized.PUT("/updatePassword", controllers.UpdatePassword)

	// Admin Management
	admin := authorized.Group("/users")
	{
		admin.POST("/add", controllers.AddAdmin)
		admin.GET("/:username", controllers.GetAdmin)
		admin.GET("/", controllers.GetAllAdmins)
		admin.PUT("/update", controllers.UpdateAdmin)
		admin.DELETE("/delete/:username", controllers.DeleteAdmin)
		admin.PUT("/inactive/:username", controllers.InactiveAdmin)
		admin.PUT("/active/:username", controllers.ActiveAdmin)
	}

	// Application Management
	r.GET("/application", controllers.GetAllApplications)
	application := authorized.Group("/application")
	{
		application.POST("/add", controllers.AddApplication)
		application.PUT("/:application_id", controllers.UpdateApplication)
		application.DELETE("/:application_id", controllers.DeleteApplication)

		// User to Application Assignment
		application.POST("/assign", controllers.AddUserToApplication)
		application.PUT("/assign/:assignment_id", controllers.UpdateUserToApplication)
		application.GET("/assignments", controllers.GetAllUserToApplications)
		application.DELETE("/assign/:assignment_id", controllers.DeleteUserToApplication)
		application.GET("/:application_id", controllers.GetApplication)

	}

	// Configuration Management
	config := authorized.Group("/config")
	{
		// config.POST("/add", controllers.CreateConfig)
		config.PUT("/update/listening-port", controllers.UpdateListeningPort)
		config.PUT("/update/rate-limit/:application_id", controllers.UpdateRateLimit)
		config.PUT("/update/remote-log-server", controllers.UpdateRemoteLogServer)
		config.PUT("/update/detect-bot/:application_id", controllers.UpdateDetectBot)
		config.PUT("/update/post-data-size/:application_id", controllers.UpdateMaxPosyDataSize)
		// config.POST("/add-app-config", controllers.CreateAppConfig)

	}
	r.GET("/config/get-app-config/:application_id", controllers.GetAppConfig)
	r.GET("/config", controllers.GetConfig)

	// Rule Management
	rules := authorized.Group("/rule")
	{
		rules.POST("/add", controllers.AddRule)
		rules.PUT("/update/:rule_id", controllers.UpdateRule)
		rules.DELETE("/delete/:rule_id", controllers.DeleteRule)
	}
	r.GET("/rule/:application_id", controllers.GetRules)
	r.GET("/rule/metadata", controllers.GetRuleMetadata)

	// Request Management with WebSocket Support
	requests := authorized.Group("/requests")
	{
		requests.GET("/", controllers.GetRequests) // Existing HTTP route to get requests
		// requests.POST("/add", controllers.AddRequest) // Existing HTTP route to add requests
		requests.GET("/stats", controllers.GetRequestStats)
		requests.GET("/blocked-stats", controllers.GetBlockedStats)
		requests.GET("/requests-per-minute", controllers.GetRequestsPerMinute)
		requests.GET("/top-blocked-countries", controllers.GetTopBlockedCountries)
		requests.GET("/all-blocked-countries", controllers.GetAllBlockedCountries)
		requests.GET("/os-stats", controllers.GetClientOSStats)
		requests.GET("/response-status-stats", controllers.GetResponseStatusStats)
		requests.GET("/most-targeted-endpoints", controllers.GetMostTargetedEndpoints)
		requests.GET("/top-attack-types", controllers.GetTopThreatTypes)
		requests.DELETE("/delete", controllers.DeleteFilteredRequests)

	}
	r.POST("/batch", controllers.HandleBatchRequests)

	// Notification Management
	notifications := authorized.Group("/notifications")
	{
		notifications.PUT("/update/:notification_id", controllers.UpdateNotification)
		notifications.GET("/all/:user_id", controllers.GetNotifications)
		notifications.DELETE("/delete/:notification_id", controllers.DeleteNotification)
	}

	// Certificate Management
	certs := authorized.Group("/certs")
	{
		certs.POST("/:application_id", controllers.AddCert) // Add a certificate
		// certs.GET("/", controllers.GetCert)               // Get certificate/key file (expects application_id & type=cert/key)
		certs.PUT("/:application_id", controllers.UpdateCert)    // Update an existing certificate
		certs.DELETE("/:application_id", controllers.DeleteCert) // Delete a certificate
	}
	r.GET("/certs", controllers.GetCert)

	// interceptor := authorized.Group("/interceptor")
	// {
	// 	// Docker Management
	// 	interceptor.GET("/start", controllers.StartInterceptor)
	// 	interceptor.GET("/stop", controllers.StopInterceptor)
	// 	interceptor.GET("/restart", controllers.ScaleInterceptor)
	// 	interceptor.GET("/repull", controllers.ScaleInterceptor)
	// }

	headers := authorized.Group("/security-headers")
	{
		headers.POST("/", controllers.AddSecurityHeader)
		headers.PUT("/:header_id", controllers.UpdateSecurityHeader)
		headers.DELETE("/:header_id", controllers.DeleteSecurityHeader)
	}
	r.GET("/security-headers/:application_id", controllers.GetSecurityHeaders)

	generateCsv := authorized.Group("/generate-csv")
	{
		generateCsv.GET("/", controllers.GenerateRequestsCSV)
	}

	notification_rule := authorized.Group("/notification-rule")
	{
		notification_rule.POST("/", controllers.AddNotificationRule)
		notification_rule.GET("/", controllers.GetNotificationRules)
		notification_rule.PUT("/:rule_id", controllers.UpdateNotificationRule)
		notification_rule.DELETE("/:rule_id", controllers.DeleteNotificationRule)
	}

	//NotificationConfig
	notification_config := authorized.Group("/notification-config")
	{
		notification_config.POST("/", controllers.AddNotificationConfig)
		notification_config.GET("/:user_id", controllers.GetNotificationConfig)
		notification_config.PUT("/:user_id", controllers.UpdateNotificationConfig)
		notification_config.DELETE("/:user_id", controllers.DeleteNotificationConfig)
	}

	r.GET("/change", controllers.CheckForUpdate)

}
