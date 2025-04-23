package models

type Request struct {
	RequestID       string  `json:"request_id" gorm:"primaryKey"`
	ApplicationName string  `json:"application_name" gorm:"not null"`
	ClientIP        string  `json:"client_ip" gorm:"not null"`
	RequestMethod   string  `json:"request_method" gorm:"not null"`
	RequestURL      string  `json:"request_url" gorm:"not null"`
	Headers         string  `json:"headers" gorm:"not null"`
	Body            string  `json:"body" gorm:"not null"`
	Timestamp       float64 `json:"timestamp" gorm:"not null"`
	ResponseCode    int     `json:"response_code" gorm:"not null"`
	Status          string  `json:"status" gorm:"not null"`
	ThreatDetected  bool    `json:"threat_detected" gorm:"not null"`
	ThreatType      string  `json:"threat_type" gorm:"not null"`
	BotDetected     bool    `json:"bot_detected" gorm:"not null"`
	GeoLocation     string  `json:"geo_location" gorm:"not null"`
	RateLimited     bool    `json:"rate_limited" gorm:"not null"`
	UserAgent       string  `json:"user_agent" gorm:"not null"`
}
