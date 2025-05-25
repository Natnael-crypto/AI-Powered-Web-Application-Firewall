package utils

import (
	"fmt"
	"os"
)

var MLEndpoint string

func InitMlService() error {
	mlHostUrl := os.Getenv("MLHOSTURL")
	if mlHostUrl == "" {
		return fmt.Errorf("MLHOST environment variable is not set")
	}
	MLEndpoint = fmt.Sprintf(mlHostUrl + "/")
	return nil
}
