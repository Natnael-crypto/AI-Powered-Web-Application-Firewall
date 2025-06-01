package main

import (
	"backend/internal/background"
	"backend/internal/config"
	"backend/internal/routes"
	"backend/internal/utils"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	config.LoadConfig()
	utils.CreateModel()
	utils.AddNotificationRule()
	err := utils.LoadIPRanges("./internal/static/iptogeo.csv")
	if err != nil {
		fmt.Println("Error loading IP ranges:", err)
		return
	}

	background.StartNotificationWatcher()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	routes.InitializeRoutes(r)

	// Start HTTPS server on port 8443
	log.Println("Starting HTTPS server on :8443")
	if err := r.RunTLS(":8443", "./certs/cert.pem", "./certs/key.pem"); err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
