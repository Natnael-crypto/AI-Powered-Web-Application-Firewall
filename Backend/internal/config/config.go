package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// JWTSecretKey holds the loaded secret key
var JWTSecretKey string
var WsKey string
var Change bool

// LoadConfig initializes environment variables
func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	WsKey = os.Getenv("WSKEY")

	if WsKey == "" {
		log.Fatalf("Missing WsKey in environment variables")
	}
	// Load JWT Secret Key
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JWTSecretKey == "" {
		log.Fatalf("Missing JWT_SECRET_KEY in environment variables")
	}

	Change = true
}
