package ml

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func EvaluateML(requestData map[string]string) (bool, string) {
	mlServerURL := "http://127.0.0.1:5001/predict"

	jsonData, _ := json.Marshal(requestData)
	resp, err := http.Post(mlServerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "Failed to connect to ML server"
	}
	defer resp.Body.Close()

	var result struct {
		Block bool   `json:"block"`
		Reason string `json:"reason"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	return result.Block, result.Reason
}
