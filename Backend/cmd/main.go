package main

import (
	"backend/internal/config"
	"backend/internal/routes"
	"backend/internal/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()
	config.LoadConfig()
	err := utils.LoadIPRanges("./internal/static/iptogeo.csv")
	if err != nil {
		fmt.Println("Error loading IP ranges:", err)
		return
	}

	// Initialize the Gin router
	r := gin.Default()

	// Initialize the routes
	routes.InitializeRoutes(r)

	// Start the server
	r.Run() // Default is localhost:8080
}
