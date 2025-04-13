package background

import (
	"backend/internal/config"
	"backend/internal/models"
	"testing"
	"time"
)

func TestProcessNotificationRules(t *testing.T) {
	// Setup test database
	config.InitDB()
	
	// Create test notification rule
	rule := models.NotificationRule{
		Name:       "Test Rule",
		HostName:   "test.com",
		Threshold:  5,
		TimeWindow: 10,
		IsActive:   true,
		UsersID:    []uint{1},
	}
	
	if err := config.DB.Create(&rule).Error; err != nil {
		t.Fatalf("Failed to create test rule: %v", err)
	}
	
	// Create test request logs
	for i := 0; i < 6; i++ {
		requestLog := models.RequestLog{
			HostName:  "test.com",
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&requestLog).Error; err != nil {
			t.Fatalf("Failed to create test request log: %v", err)
		}
	}
	
	// Run the notification processor
	processNotificationRules()
	
	// Verify notification was created
	var notification models.Notification
	if err := config.DB.First(&notification).Error; err != nil {
		t.Errorf("Expected notification to be created, but got error: %v", err)
	}
	
	// Cleanup
	config.DB.Unscoped().Delete(&rule)
	config.DB.Unscoped().Delete(&models.RequestLog{})
	config.DB.Unscoped().Delete(&notification)
}

func TestShouldTriggerNotification(t *testing.T) {
	// Setup test database
	config.InitDB()
	
	// Create test rule
	rule := models.NotificationRule{
		Name:       "Test Rule",
		HostName:   "test.com",
		Threshold:  3,
		TimeWindow: 5,
		IsActive:   true,
	}
	
	// Test case 1: Below threshold
	for i := 0; i < 2; i++ {
		requestLog := models.RequestLog{
			HostName:  "test.com",
			CreatedAt: time.Now(),
		}
		if err := config.DB.Create(&requestLog).Error; err != nil {
			t.Fatalf("Failed to create test request log: %v", err)
		}
	}
	
	if shouldTriggerNotification(rule) {
		t.Error("Expected notification not to trigger when below threshold")
	}
	
	// Test case 2: At threshold
	requestLog := models.RequestLog{
		HostName:  "test.com",
		CreatedAt: time.Now(),
	}
	if err := config.DB.Create(&requestLog).Error; err != nil {
		t.Fatalf("Failed to create test request log: %v", err)
	}
	
	if !shouldTriggerNotification(rule) {
		t.Error("Expected notification to trigger when at threshold")
	}
	
	// Cleanup
	config.DB.Unscoped().Delete(&models.RequestLog{})
}

func TestCreateNotification(t *testing.T) {
	// Setup test database
	config.InitDB()
	
	// Create test rule
	rule := models.NotificationRule{
		Name: "Test Rule",
	}
	
	// Test notification creation
	createNotification(rule)
	
	// Verify notification was created
	var notification models.Notification
	if err := config.DB.First(&notification).Error; err != nil {
		t.Errorf("Expected notification to be created, but got error: %v", err)
	}
	
	// Cleanup
	config.DB.Unscoped().Delete(&notification)
}

func TestSendEmail(t *testing.T) {
	// Setup test database
	config.InitDB()
	
	// Create test users
	users := []models.NotificationConfig{
		{
			UserID: 1,
			Email:   "test1@example.com",
		},
		{
			UserID: 2,
			Email:   "test2@example.com",
		},
	}
	
	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
	}
	
	// Create test rule
	rule := models.NotificationRule{
		Name:    "Test Rule",
		UsersID: []uint{1, 2},
	}
	
	// Test email sending
	sendEmail(rule)
	
	// Cleanup
	config.DB.Unscoped().Delete(&models.NotificationConfig{})
} 