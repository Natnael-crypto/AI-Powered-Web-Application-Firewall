package config

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	RedirectURL string `yaml:"redirect_url"`
	Port        int    `yaml:"port"`
	DebugMode   bool   `yaml:"debug_mode"`
}

type AppConfig struct {
	Server ServerConfig `yaml:"server"`
}

var Config AppConfig

func ParseConfig() {
	// Attempt to load from environment variables
	if port, exists := os.LookupEnv("SERVER_PORT"); exists {
		parsedPort, err := strconv.Atoi(port)
		if err == nil {
			Config.Server.Port = parsedPort
		}
	}

	// Load from config file
	file, err := os.Open("./internal/config/config.yaml")
	if err != nil {
		log.Println("Config file not found, skipping file load.")
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	// Fallback validation
	if Config.Server.Port == 0 {
		Config.Server.Port = 8080 // Default port
	}
}
