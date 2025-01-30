package config

import (
	"backend/internal/models" // Adjust the import path to where your models are located
	"fmt"
	"log"
	"os"

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
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Database migration completed successfully.")

	// Create indexes to optimize queries
	createIndexes()
}

// createIndexes ensures database indexes are created for performance optimization
func createIndexes() {
	log.Println("Creating necessary indexes...")

	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_requests_application_name ON requests(application_name);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_client_ip ON requests(client_ip);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_request_method ON requests(request_method);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_timestamp ON requests(timestamp);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_threat_detected ON requests(threat_detected);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_bot_detected ON requests(bot_detected);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_rate_limited ON requests(rate_limited);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_threat_type ON requests(threat_type);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_geo_location ON requests(geo_location);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_user_agent ON requests(user_agent);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_action_taken ON requests(action_taken);`,
		`CREATE INDEX IF NOT EXISTS idx_requests_full_text ON requests USING GIN (to_tsvector('english', headers || ' ' || body || ' ' || request_url));`,
	}

	for _, query := range indexQueries {
		if err := DB.Exec(query).Error; err != nil {
			log.Printf("Failed to create index: %v", err)
		}
	}

	log.Println("Indexes created successfully.")
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
