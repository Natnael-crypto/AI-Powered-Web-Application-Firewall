package models

import (
	"time"
)

type Notification struct {
	NotificationID string    `json:"notification_id" gorm:"primaryKey"`
	UserID         string    `json:"user_id" gorm:"not null" `
	Message        string    `json:"message" gorm:"not null" `
	Timestamp      time.Time `json:"timestamp"`
	Status         bool      `json:"status" gorm:"not null" `
}

type NotificationRule struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name" gorm:"not null" `
	ThreatType string    `json:"threat_type" gorm:"unique;not null" `
	Threshold  int       `json:"threshold" gorm:"not null" `
	TimeWindow int       `json:"time_window" gorm:"not null" `
	IsActive   bool      `json:"is_active" gorm:"not null" `
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type NotificationRuleToApplication struct {
	ID                 string `json:"id" gorm:"primaryKey" `
	NotificationRuleID string `json:"notification_rule_id" gorm:"not null" `
	ApplicationID      string `json:"application_id" gorm:"not null" `
}

type NotificationRuleToUser struct {
	ID                 string `json:"id" gorm:"primaryKey" `
	NotificationRuleID string `json:"notification_rule_id" gorm:"not null" `
	UserID             string `json:"application_id" gorm:"not null" `
}

type NotificationConfig struct {
	UserID string `json:"user_id" gorm:"primaryKey"`
	Email  string `json:"email"`
}

type NotificationSender struct {
	Email       string `json:"sender_email" gorm:"primaryKey" binding:"required,email"`
	AppPassword string `json:"app_password" gorm:"not null" binding:"required"`
}

type NotificationInput struct {
	UserID           string `json:"user_id" `
	NotificationType string `json:"notification_type" `
	Message          string `json:"message" `
	Status           bool   `json:"status" `
	Severity         string `json:"severity" `
}
