package config

import (
	"backend/internal/models" // Adjust the import path to where your models are located
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection using GORM
func InitDB() {
	// Database connection configuration
	const (
		host     = "pg-2c27c868-qzueos-e68f.e.aivencloud.com" // Update as needed
		port     = 25211                                      // Default PostgreSQL port
		user     = "avnadmin"                                 // Replace with your PostgreSQL username
		password = "AVNS_rcBeTPKpktP4be5ZEZx"                 // Replace with your PostgreSQL password
		dbname   = "waf"                                      // Replace with your database name
	)

	// Construct the connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Open the database connection using GORM
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify the connection
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL using GORM!")

	// Run migrations to create the tables if they don't exist
	err = DB.AutoMigrate(&models.Application{}, &models.User{}, models.UserToApplication{}, models.Conf{}, models.Rule{}, models.Request{}, models.Notification{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migration completed successfully.")
}

// CloseDB closes the GORM database connection
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting SQL DB object: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
