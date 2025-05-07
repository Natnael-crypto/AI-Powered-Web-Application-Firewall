package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	backendHost := os.Getenv("BACKENDHOST")
	if backendHost == "" {
		return nil, fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		return nil, fmt.Errorf("BACKENDPORT environment variable is not set")
	}

	url := fmt.Sprintf("http://%s:%s/interceptor/rule/%s", backendHost, backendPort, applicationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Service", "I")

	client := &http.Client{}
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
	fileName := fmt.Sprintf("%s.conf", applicationID)

	file, err := os.Create("./internal/config/custom/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	for _, rule := range rules {
		_, err := file.WriteString(rule.RuleString + "\n")
		if err != nil {
			return "", fmt.Errorf("error writing to file: %v", err)
		}
	}

	return fileName, nil
}
