package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// fetchCert fetches a certificate and key for the given application ID,
// writes them to files in ./certs/, and returns the file paths.
func fetchCert(applicationID string) (string, string, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDHOST")
	backendPort := os.Getenv("BACKENDPORT")

	if backendHost == "" || backendPort == "" {
		return "", "", fmt.Errorf("BACKENDHOST or BACKENDPORT environment variable is not set")
	}

	// Construct URLs to fetch cert and key
	backendURLCert := fmt.Sprintf("http://%s:%s/certs?application_id=%s&type=cert", backendHost, backendPort, applicationID)
	backendURLKey := fmt.Sprintf("http://%s:%s/certs?application_id=%s&type=key", backendHost, backendPort, applicationID)

	// Fetch certificate
	certPath, err := fetchAndSaveFile(backendURLCert, applicationID, "crt")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch and save cert: %v", err)
	}

	// Fetch key
	keyPath, err := fetchAndSaveFile(backendURLKey, applicationID, "key")
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch and save key: %v", err)
	}

	// Return paths of saved files
	return certPath, keyPath, nil
}

// fetchAndSaveFile downloads the file from the given URL and saves it to ./certs/<applicationID>.<ext>
func fetchAndSaveFile(url, applicationID, ext string) (string, error) {
	// Send GET request to fetch the file
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	// Create certs directory if it doesn't exist
	certDir := "certs"
	if err := os.MkdirAll(certDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create certs directory: %v", err)
	}

	// Define file path
	filePath := filepath.Join(certDir, fmt.Sprintf("%s.%s", applicationID, ext))

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	return filePath, nil
}
