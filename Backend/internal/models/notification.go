package models

import "time"

type Notification struct {
	NotificationID   string    `json:"notification_id"`
	UserID           string    `json:"user_id"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	Timestamp        time.Time `json:"timestamp"`
	Status           bool      `json:"status"`
	Severity         string    `json:"severity"`
}
