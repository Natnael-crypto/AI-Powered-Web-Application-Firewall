package routes

import (
	"backend/internal/controllers"
	middleware "backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired)

	allowed_ip := r.Group("/")
	allowed_ip.Use(middleware.AllowlistMiddleware)

	authorized.PUT("/updatePassword", controllers.UpdatePassword)
	authorized.GET("/is-logged-in", controllers.IsLoggedIN)

	admin := authorized.Group("/users")
	{
		admin.POST("/add", controllers.AddAdmin)
		admin.GET("/:username", controllers.GetAdmin)
		admin.GET("/id/:user_id", controllers.GetAdminByID)
		admin.GET("/", controllers.GetAllAdmins)
		admin.DELETE("/delete/:username", controllers.DeleteAdmin)
		admin.PUT("/inactive/:username", controllers.InactiveAdmin)
		admin.PUT("/active/:username", controllers.ActiveAdmin)
	}

	application := authorized.Group("/application")
	{
		application.POST("/add", controllers.AddApplication)
		application.PUT("/:application_id", controllers.UpdateApplication)
		application.DELETE("/:application_id", controllers.DeleteApplication)
		application.GET("", controllers.GetAllApplicationsAdmin)
		application.POST("/assign", controllers.AddUserToApplication)
		application.PUT("/assign/:assignment_id", controllers.UpdateUserToApplication)
		application.GET("/assignments", controllers.GetAllUserToApplications)
		application.DELETE("/assign/:assignment_id", controllers.DeleteUserToApplication)
		application.GET("/:application_id", controllers.GetApplication)
	}

	config := authorized.Group("/config")
	{
		config.PUT("/update/listening-port", controllers.UpdateListeningPort)
		config.PUT("/update/rate-limit/:application_id", controllers.UpdateRateLimit)
		config.PUT("/update/remote-log-server", controllers.UpdateRemoteLogServer)
		config.PUT("/update/detect-bot/:application_id", controllers.UpdateDetectBot)
		config.PUT("/update/post-data-size/:application_id", controllers.UpdateMaxPosyDataSize)
		config.PUT("/update/tls/:application_id", controllers.UpdateTls)
		config.GET("/", controllers.GetConfigAdmin)
		config.GET("/:application_id", controllers.GetAppConfig)
	}

	sysEmail := authorized.Group("/sys-email")
	{
		sysEmail.POST("/", controllers.AddEmail)
		sysEmail.GET("/", controllers.GetEmail)
		sysEmail.PUT("/", controllers.UpdateEmail)
	}

	rules := authorized.Group("/rule")
	{
		rules.POST("/add", controllers.AddRule)
		rules.PUT("/update/:rule_id", controllers.UpdateRule)
		rules.DELETE("/delete/:rule_id", controllers.DeleteRule)
		rules.POST("/deactivate/:rule_id", controllers.DeactivateRule)
		rules.POST("/activate/:rule_id", controllers.ActivateRule)
		rules.GET("", controllers.GetAllRulesAdmin)
		rules.GET("/:rule_id", controllers.GetOneRule)
	}

	requests := authorized.Group("/requests")
	{
		requests.GET("", controllers.GetRequests)
		requests.GET("/overall-stat", controllers.GetOverallStats)
		requests.GET("/:request_id", controllers.GetRequestByID)
		requests.GET("/requests-per-minute", controllers.GetRequestRateLastMinute)
		requests.GET("/all-countries-stat", controllers.GetAllCountriesStat)
		requests.GET("/os-stats", controllers.GetClientOSStats)
		requests.GET("/response-status-stat", controllers.GetResponseStatusStats)
		requests.GET("/most-targeted-endpoints", controllers.GetMostTargetedEndpoints)
		requests.GET("/top-attack-types", controllers.GetTopThreatTypes)
		requests.DELETE("/delete", controllers.DeleteFilteredRequests)
	}

	notifications := authorized.Group("/notifications")
	{
		notifications.PUT("/update/:notification_id", controllers.UpdateNotification)
		notifications.GET("/all/:user_id", controllers.GetNotifications)
		notifications.DELETE("/delete/:notification_id", controllers.DeleteNotification)
	}

	certs := authorized.Group("/certs")
	{
		certs.GET("/certs", controllers.GetCertAdmin)
		certs.POST("/:application_id", controllers.AddCert)
		certs.PUT("/:application_id", controllers.UpdateCert)
		certs.DELETE("/:application_id", controllers.DeleteCert)
	}

	interceptor := authorized.Group("/interceptor")
	{
		interceptor.GET("/start", controllers.StartInterceptor)
		interceptor.GET("/stop", controllers.StopInterceptor)
		interceptor.GET("/restart", controllers.RestartInterceptor)
	}

	headers := authorized.Group("/security-headers")
	{
		headers.POST("", controllers.AddSecurityHeader)
		headers.PUT("/:header_id", controllers.UpdateSecurityHeader)
		headers.DELETE("/:header_id", controllers.DeleteSecurityHeader)
		headers.GET("", controllers.GetSecurityHeadersAdmin)

	}

	generateCsv := authorized.Group("/generate-csv")
	{
		generateCsv.GET("", controllers.GenerateRequestsCSV)
	}

	notification_rule := authorized.Group("/notification-rule")
	{
		notification_rule.POST("/", controllers.AddNotificationRule)
		notification_rule.GET("/:application_id", controllers.GetNotificationRule)
		notification_rule.GET("", controllers.GetNotificationRules)
		notification_rule.PUT("/:rule_id", controllers.UpdateNotificationRule)
		notification_rule.DELETE("/:rule_id", controllers.DeleteNotificationRule)
	}

	notification_config := authorized.Group("/notification-config")
	{
		notification_config.POST("", controllers.AddNotificationConfig)
		notification_config.GET("", controllers.GetNotificationConfig)
		notification_config.GET("/all", controllers.GetAllNotificationConfig)
		notification_config.PUT("/:user_id", controllers.UpdateNotificationConfig)
		notification_config.DELETE("/:user_id", controllers.DeleteNotificationConfig)
	}

	ai_analysis := authorized.Group("/")
	{
		ai_analysis.GET("/model/select/:model_id", controllers.SelectActiveModel)
		ai_analysis.GET("/models", controllers.GetModels)
		ai_analysis.PUT("/model/update/setting", controllers.UpdateTrainingSettings)
	}

	ms_services := allowed_ip.Group("/ml")
	{
		ms_services.POST("/submit-analysis", controllers.SubmitAnalysisResults)
		ms_services.GET("/model/untrained", controllers.GetModelForMLs)
		ms_services.POST("/model/results", controllers.SubmitTrainResults)
		ms_services.GET("/model/selected", controllers.GetSelectedModel)
		ms_services.GET("/changes", controllers.MlCheckState)
	}

	interceptor_services := allowed_ip.Group("/interceptor")
	{
		interceptor_services.GET("/application", controllers.GetAllApplications)
		interceptor_services.POST("/batch", controllers.HandleBatchRequests)
		interceptor_services.GET("/config", controllers.GetConfig)
		interceptor_services.GET("/rule/:application_id", controllers.GetRules)
		interceptor_services.GET("/certs", controllers.GetCert)
		interceptor_services.GET("/is-running", controllers.InterceptorCheckState)
		interceptor_services.GET("/security-headers/:application_id", controllers.GetSecurityHeaders)
	}

	services := authorized.Group("/service")
	{
		services.GET("", controllers.GetAllowedIps)
		services.POST("", controllers.AddAllowedIp)
		services.PUT("/:id", controllers.UpdateAllowedIp)
		services.DELETE("/:id", controllers.DeleteAllowedIp)
	}

}
