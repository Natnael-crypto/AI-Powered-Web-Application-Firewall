package config

import (
	"backend/internal/models" // Adjust the import path to where your models are located
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection using GORM
func InitDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	// Read database credentials from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Ensure required variables are set
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("Missing required database environment variables")
	}

	// Construct the connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

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
	err = DB.AutoMigrate(
		&models.Application{},
		&models.User{},
		&models.UserToApplication{},
		&models.Conf{},
		&models.Rule{},
		&models.Request{},
		&models.Notification{},
		&models.Cert{},
		&models.AppConf{},
		&models.NotificationRule{},
		&models.NotificationConfig{},
		DB.AutoMigrate(&models.SecurityHeader{}),
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migration completed successfully.")

	// Create indexes to optimize queries

	newConf := models.Conf{
		ID:              uuid.New().String(),
		ListeningPort:   "80",
		RemoteLogServer: "",
	}

	if err := CreateConfigLocal(newConf); err != nil {
		fmt.Println("Unable to set Default Listening Port 80")
	}
}

func CreateConfigLocal(conf models.Conf) error {
	var existingConfig models.Conf
	if err := DB.First(&existingConfig).Error; err == nil {
		return nil
	}
	if err := DB.Create(&conf).Error; err != nil {
		return err
	}
	return nil
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
