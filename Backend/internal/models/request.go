package models

import "time"

// Request represents a log of a web request
type Request struct {
	RequestID        string    `json:"request_id" gorm:"primaryKey"`
	ApplicationID    string    `json:"application_id"`
	ClientIP         string    `json:"client_ip"`
	RequestMethod    string    `json:"request_method"`
	RequestURL       string    `json:"request_url"`
	Headers          string    `json:"headers"`
	Body             string    `json:"body"`
	Timestamp        time.Time `json:"timestamp"`
	ResponseCode     int       `json:"response_code"`
	Status           string    `json:"status"`
	MatchedRules     string    `json:"matched_rules"`
	ThreatDetected   bool      `json:"threat_detected"`
	ThreatType       string    `json:"threat_type"`
	ActionTaken      string    `json:"action_taken"`
	BotDetected      bool      `json:"bot_detected"`
	GeoLocation      string    `json:"geo_location"`
	RateLimited      bool      `json:"rate_limited"`
	UserAgent        string    `json:"user_agent"`
	AIAnalysisResult string    `json:"ai_analysis_result"`
}
