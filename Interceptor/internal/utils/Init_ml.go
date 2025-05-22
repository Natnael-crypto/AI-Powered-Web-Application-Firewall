package utils

import (
	"fmt"
	"os"
)

var MLEndpoint string

func InitMlService() error {
	mlHost := os.Getenv("MLHOST")
	if mlHost == "" {
		return fmt.Errorf("MLHOST environment variable is not set")
	}

	mlPort := os.Getenv("MLPORT")
	if mlPort == "" {
		return fmt.Errorf("MLPORT environment variable is not set")
	}
	MLEndpoint = fmt.Sprintf("http://%s:%s/", mlHost, mlPort)
	return nil
}
