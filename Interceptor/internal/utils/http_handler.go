package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type MessageModel struct {
	RequestID       string `json:"request_id"`
	ApplicationName string `json:"application_name"`
	ClientIP        string `json:"client_ip"`
	RequestMethod   string `json:"request_method"`
	RequestURL      string `json:"request_url"`
	Headers         string `json:"headers"`
	Body            string `json:"body"`
	ResponseCode    int    `json:"response_code"`
	Status          string `json:"status"`
	ThreatDetected  bool   `json:"threat_detected"`
	ThreatType      string `json:"threat_type"`
	BotDetected     bool   `json:"bot_detected"`
	GeoLocation     string `json:"geo_location"`
	RateLimited     bool   `json:"rate_limited"`
	UserAgent       string `json:"user_agent"`
	Token           string `json:"token"`
	AIResult        bool   `json:"ai_result"`
	AIThreatType    string `json:"ai_threat_type"`
	RuleDetected    bool   `json:"rule_detected"`
}

var (
	messageQueue    []MessageModel
	queueMutex      sync.Mutex
	queueLimit      = 1000
	sendInterval    = 60 * time.Second
	httpClient      = &http.Client{Timeout: 10 * time.Second}
	backendEndpoint string
)

func InitHttpHandler() error {
	backendHost := os.Getenv("BACKENDURL")

	backendEndpoint = fmt.Sprintf(backendHost + "/interceptor/batch")

	go StartBatchSender()
	return nil
}

func SendToBackend(message MessageModel) {
	queueMutex.Lock()
	messageQueue = append(messageQueue, message)
	queueFull := len(messageQueue) >= queueLimit
	queueMutex.Unlock()

	if queueFull {
		go flushQueue()
	}
}

func flushQueue() {
	queueMutex.Lock()
	batch := messageQueue
	messageQueue = nil
	queueMutex.Unlock()

	if len(batch) == 0 {
		return
	}

	data, err := json.Marshal(batch)
	if err != nil {
		log.Printf("Failed to marshal batch: %v", err)
		return
	}

	req, err := http.NewRequest("POST", backendEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "I")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to send batch to backend: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Batch sent to backend, status: %s", resp.Status)
}

func StartBatchSender() {
	ticker := time.NewTicker(sendInterval)
	defer ticker.Stop()

	for range ticker.C {
		go flushQueue()
	}
}
