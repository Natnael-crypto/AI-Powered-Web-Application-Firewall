package ml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interceptor/internal/utils"
	"net/http"
)

type RequestData struct {
	RequestID string `json:"request_id"`
	Url       string `json:"url"`
	Headers   string `json:"headers"`
	Body      string `json:"body"`
}

func EvaluateML(requestData RequestData) (bool, float64, error) {
	mlServerURL := utils.MLEndpoint + "predict"

	jsonData, _ := json.Marshal(requestData)

	req, err := http.NewRequest("POST", mlServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, 0.0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Service", "I")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return false, 0.0, err
	}

	defer resp.Body.Close()

	var result struct {
		Block   bool    `json:"block"`
		Percent float64 `json:"percent"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	return result.Block, result.Percent, nil
}
