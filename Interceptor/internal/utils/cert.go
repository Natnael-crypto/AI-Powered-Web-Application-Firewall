package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)


func FetchCert(applicationID string) (string, string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDHOST")
	backendPort := os.Getenv("BACKENDPORT")

	if backendHost == "" || backendPort == "" {
		return "", "", fmt.Errorf("BACKENDHOST or BACKENDPORT environment variable is not set")
	}

	backendURLCert := fmt.Sprintf("http://%s:%s/certs?application_id=%s&type=cert", backendHost, backendPort, applicationID)
	backendURLKey := fmt.Sprintf("http://%s:%s/certs?application_id=%s&type=key", backendHost, backendPort, applicationID)

	certPath, err := fetchAndSaveFile(backendURLCert, applicationID, "crt")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch and save cert: %v", err)
	}

	keyPath, err := fetchAndSaveFile(backendURLKey, applicationID, "key")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch and save key: %v", err)
	}

	return certPath, keyPath, nil
}

func fetchAndSaveFile(url, applicationID, ext string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	certDir := "certs"
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create certs directory: %v", err)
	}

	filePath := filepath.Join(certDir, fmt.Sprintf("%s.%s", applicationID, ext))

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	return filePath, nil
}
