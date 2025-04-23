package models

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	NotificationID string    `json:"notification_id" gorm:"primaryKey"`
	UserID         string    `json:"user_id" gorm:"not null" `
	Message        string    `json:"message" gorm:"not null" `
	Timestamp      time.Time `json:"timestamp"`
	Status         bool      `json:"status" gorm:"not null" `
}

type NotificationRule struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	CreatedBy  string         `json:"created_by" gorm:"not null" `
	Name       string         `json:"name" gorm:"not null" `
	HostName   string         `json:"hostname" gorm:"not null"`
	ThreatType string         `json:"threat_type" gorm:"not null" `
	Threshold  int            `json:"threshold" gorm:"not null" `
	TimeWindow int            `json:"time_window" gorm:"not null" `
	IsActive   bool           `json:"is_active" gorm:"not null" `
	UsersID    datatypes.JSON `gorm:"type:jsonb" json:"users_id" `
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
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
