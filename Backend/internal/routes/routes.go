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
	admin := authorized.Group("/admin")
	{
		admin.POST("/add", controllers.AddAdmin)
		admin.GET("/:username", controllers.GetAdmin)
		admin.GET("/", controllers.GetAllAdmins)
		admin.PUT("/update", controllers.UpdateAdmin)
		admin.DELETE("/delete/:username", controllers.DeleteAdmin)
	}

	// Application Management
	r.GET("/application", controllers.GetAllApplications)
	r.GET("/:application_id", controllers.GetApplication)
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
	}

	// Configuration Management
	config := authorized.Group("/config")
	{
		config.POST("/add", controllers.CreateConfig)
		config.PUT("/update/:id", controllers.UpdateConfig)
	}
	r.GET("/config", controllers.GetConfig)

	// Rule Management
	rules := authorized.Group("/rule")
	{
		rules.POST("/add", controllers.AddRule)
		rules.PUT("/update/:rule_id", controllers.UpdateRule)
		rules.DELETE("/delete/:rule_id", controllers.DeleteRule)
	}
	r.GET("/rule/:application_id", controllers.GetRules)

	// Request Management with WebSocket Support
	requests := authorized.Group("/requests")
	{
		requests.GET("/", controllers.GetRequests)    // Existing HTTP route to get requests
		requests.POST("/add", controllers.AddRequest) // Existing HTTP route to add requests
	}
	r.GET("/ws", controllers.HandleWebSocket)

	// Notification Management
	notifications := authorized.Group("/notifications")
	{
		notifications.POST("/add", controllers.CreateNotification)
		notifications.PUT("/update/:notification_id", controllers.UpdateNotification)
		notifications.GET("/all/:user_id", controllers.GetNotifications)
		notifications.DELETE("/delete/:notification_id", controllers.DeleteNotification)
	}

	// Certificate Management
	certs := authorized.Group("/certs")
	{
		certs.POST("/", controllers.AddCert) // Add a certificate
		// certs.GET("/", controllers.GetCert)               // Get certificate/key file (expects application_id & type=cert/key)
		certs.PUT("/:application_id", controllers.UpdateCert)    // Update an existing certificate
		certs.DELETE("/:application_id", controllers.DeleteCert) // Delete a certificate
	}
	r.GET("/certs", controllers.GetCert)

	docker := authorized.Group("/docker")
	{
		// Docker Management
		docker.GET("/start", controllers.StartContainer)
		docker.GET("/stop", controllers.StopContainer)
		docker.GET("/restart", controllers.RestartContainer)
	}

}
