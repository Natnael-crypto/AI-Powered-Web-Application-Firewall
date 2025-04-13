package models

import "time"

type Notification struct {
	NotificationID   string    `json:"notification_id" gorm:"primaryKey"`
	UserID           string    `json:"user_id" `
	NotificationType string    `json:"notification_type" `
	Message          string    `json:"message" `
	Timestamp        time.Time `json:"timestamp" `
	Status           bool      `json:"status" `
	Severity         string    `json:"severity" `
}

type NotificationRule struct {
	ID         string      `gorm:"primaryKey" json:"id"`
	CreatedBy  string		`json:"created_by" `
	Name       string    `json:"name" `
	HostName   string    `json:"hostname" `
	RuleType   string    `json:"rule_type" `
	Threshold  int       `json:"threshold" `
	TimeWindow int       `json:"time_window" `
	Severity   string    `json:"severity" `
	IsActive   bool      `json:"is_active" `
	UsersID    []string  `gorm:"type:json" json:"users_id" ` // List of users to notify
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NotificationConfig struct {
	UserID string `json:"user_id" gorm:"primaryKey"`
	Email  string `json:"email"`
}

type NotificationInput struct {
	UserID           string `json:"user_id" `
	NotificationType string `json:"notification_type" `
	Message          string `json:"message" `
	Status           bool   `json:"status" `
	Severity         string `json:"severity" `
}
