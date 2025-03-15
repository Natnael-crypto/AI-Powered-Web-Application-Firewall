package models

import "time"

type Notification struct {
	NotificationID   string    `json:"notification_id" gorm:"primaryKey"`
	UserID           string    `json:"user_id"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	Timestamp        time.Time `json:"timestamp"`
	Status           bool      `json:"status"`
	Severity         string    `json:"severity"`
}

type NotificationRule struct {
	NotificationRuleID string `json:"notification_rule_id" gorm:"primaryKey"`
	UserID             string `json:"user_id"`
	NotificationType   string `json:"notification_type"`
	Description        string `json:"description"`
	Threshold          int    `json:"threshold"`
	Active             bool   `json:"active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type NotificationType string

const (
	Critical NotificationType = "Critical"
	High NotificationType = "High"
	Medium NotificationType = "Medium"
	Low NotificationType = "Low"
)

type NotificationConfig struct {
	NotificationConfigID string `json:"notification_config_id" gorm:"primaryKey"`
	Email                string `json:"email"`
}

