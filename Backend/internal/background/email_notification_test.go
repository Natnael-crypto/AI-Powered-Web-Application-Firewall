package background_test

import (
	"backend/internal/background"
	"backend/internal/models"
	"os"
	"testing"
	"time"
)

func TestSendEmail(t *testing.T) {
	recipient := models.NotificationConfig{
		UserID: "123",
		Email:  os.Getenv("RECIPIENT_EMAIL"),
	}

	sender := models.NotificationSender{
		Email:       os.Getenv("SENDER_EMAIL"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	triggeredRule := models.NotificationRule{
		ID:         "123",
		CreatedBy:  "super admin",
		Name:       "SQL Injection",
		ThreatType: "Attack",
		Threshold:  5,
		TimeWindow: 1,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	background.SendEmail(recipient, sender, triggeredRule, "This is a test email")
}
