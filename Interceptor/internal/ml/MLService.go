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

func EvaluateML(requestData RequestData) (bool, float64, float64, error) {
	mlServerURL := utils.MLEndpoint + "predict"

	jsonData, _ := json.Marshal(requestData)

	requestData.Body = utils.RecursiveDecode(requestData.Body, 5)
	requestData.Url = utils.RecursiveDecode(requestData.Url, 5)
	requestData.Headers = utils.RecursiveDecode(requestData.Headers, 5)

	req, err := http.NewRequest("POST", mlServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, 0.0, 0.0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "I")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err)
		return false, 0.0, 0.0, fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	var result struct {
		Success    bool    `json:"success"`
		Prediction string  `json:"prediction"`
		Normal     float64 `json:"Normal"`
		Anomaly    float64 `json:"Anomaly"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	Block := false

	if result.Prediction == "Anomaly" {
		Block = true
	}

	return Block, result.Normal, result.Anomaly, nil
}
