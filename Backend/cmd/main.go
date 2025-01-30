package main

import (
	"backend/internal/config"
	"backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()

	// Initialize the Gin router
	r := gin.Default()

	// Initialize the routes
	routes.InitializeRoutes(r)

	// Start the server
	r.Run() // Default is localhost:8080
}
