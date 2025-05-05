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
	err := utils.LoadIPRanges("./internal/static/iptogeo.csv")
	if err != nil {
		fmt.Println("Error loading IP ranges:", err)
		return
	}

	background.StartNotificationWatcher()
	log.Println("Notification watcher started")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.InitializeRoutes(r)

	log.Printf("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
