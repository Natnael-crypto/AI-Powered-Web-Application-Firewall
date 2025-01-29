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
	application := authorized.Group("/application")
	{
		application.POST("/add", controllers.AddApplication)
		application.GET("/:application_id", controllers.GetApplication)
		application.GET("/", controllers.GetAllApplications)
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
		rules.GET("/all/:application_id", controllers.GetRules)
		rules.DELETE("/delete/:rule_id", controllers.DeleteRule)
	}

	// Request Management with WebSocket Support
	requests := authorized.Group("/requests")
	{
		requests.GET("/", controllers.GetRequests)    // Existing HTTP route to get requests
		requests.POST("/add", controllers.AddRequest) // Existing HTTP route to add requests
		// requests.GET("/user/ws", controllers.HandleWebSocket)       // WebSocket route for real-time updates
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
}
