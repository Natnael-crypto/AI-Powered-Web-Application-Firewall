package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Rule struct to map the rule data from the API response
type Rule struct {
	RuleID         string `json:"rule_id"`
	RuleType       string `json:"rule_type"`
	RuleMethod     string `json:"rule_method"`
	RuleDefinition string `json:"rule_definition"`
	Action         string `json:"action"`
	ApplicationID  string `json:"application_id"`
	RuleString     string `json:"rule_string"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	IsActive       bool   `json:"is_active"`
	Category       string `json:"category"`
}

// RulesResponse struct to map the API response
type RulesResponse struct {
	Rules []Rule `json:"rules"`
}

// FetchRules function to make the API call and fetch rules for a given application ID
func FetchRules(applicationID string) (*RulesResponse, error) {
	// Get the backend host and port from environment variables
	backendHost := os.Getenv("BACKENDHOST")
	if backendHost == "" {
		return nil, fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		return nil, fmt.Errorf("BACKENDPORT environment variable is not set")
	}

	// Construct the URL
	url := fmt.Sprintf("http://%s:%s/rule/%s", backendHost, backendPort, applicationID)

	// Make the GET request
	resp, err := http.Get(url) // #nosec G107 -- The url is constructed from an untainted sourcce
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the response body into the RulesResponse struct
	var rulesResponse RulesResponse
	err = json.Unmarshal(body, &rulesResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return &rulesResponse, nil
}

// WriteRuleToFile function to write the rule string to a file
func WriteRuleToFile(applicationID string, rules []Rule) (string, error) {
	// Create the filename based on application ID
	fileName := fmt.Sprintf("%s.conf", applicationID)

	// Open the file for writing
	file, err := os.Create("./internal/config/custom/" + fileName) // #nosec G304 -- The file path comes from an untainted source
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the rule strings to the file
	for _, rule := range rules {
		_, err := file.WriteString(rule.RuleString + "\n")
		if err != nil {
			return "", fmt.Errorf("error writing to file: %v", err)
		}
	}

	// Return the file name
	return fileName, nil
}
