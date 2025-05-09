package utils

import (
	"fmt"
	"os"
)

var MLEndpoint string

func InitMlService() error {
	mlHost := os.Getenv("BACKENDHOST")
	if mlHost == "" {
		return fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	mlPort := os.Getenv("BACKENDPORT")
	if mlPort == "" {
		return fmt.Errorf("BACKENDPORT environment variable is not set")
	}
	MLEndpoint = fmt.Sprintf("http://%s:%s/", mlHost, mlPort)
	return nil
}
