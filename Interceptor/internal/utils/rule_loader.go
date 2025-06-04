package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Rule struct {
	RuleID         string `json:"rule_id"`
	RuleType       string `json:"rule_type"`
	RuleMethod     string `json:"rule_method"`
	RuleDefinition string `json:"rule_definition"`
	Action         string `json:"action"`
	RuleString     string `json:"rule_string"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	IsActive       bool   `json:"is_active"`
	Category       string `json:"category"`
}

type RulesResponse struct {
	Rules []Rule `json:"rules"`
}

func FetchRules(applicationID string) (*RulesResponse, error) {
	backendHost := os.Getenv("BACKENDURL")

	url := fmt.Sprintf(backendHost+"/interceptor/rule/%s", applicationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Service", "I")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var rulesResponse RulesResponse
	err = json.Unmarshal(body, &rulesResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &rulesResponse, nil
}

func WriteRuleToFile(applicationID string, rules []Rule) (string, error) {
	dirPath := "./internal/config/custom/"
	fileName := fmt.Sprintf("%s.conf", applicationID)
	fullPath := filepath.Join(dirPath, fileName)

	// Check if folder exists
	if _, err := os.Stat(dirPath); err == nil {
		// Folder exists, remove it
		if err := os.RemoveAll(dirPath); err != nil {
			return "", fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// Recreate the folder
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Create the rule file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write rules to file
	for _, rule := range rules {
		_, err := file.WriteString(rule.RuleString + "\n")
		if err != nil {
			return "", fmt.Errorf("error writing to file: %v", err)
		}
	}

	return fileName, nil
}
