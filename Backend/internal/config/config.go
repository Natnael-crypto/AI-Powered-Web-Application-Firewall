package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	JWTSecretKey        string
	WsKey               string
	InterceptorRunning  bool = false
	Change              bool = false
	ModelSettingUpdated bool = false
	SelectModel         bool = false
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	WsKey = os.Getenv("WSKEY")

	if WsKey == "" {
		log.Fatalf("Missing WsKey in environment variables")
	}
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JWTSecretKey == "" {
		log.Fatalf("Missing JWT_SECRET_KEY in environment variables")
	}

	InterceptorRunning = false
	Change = false
}
