package background

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"log"
	"time"
)

// StartNotificationWatcher runs a background job every 5 minutes
// to check notification rules.
func StartNotificationWatcher() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			processNotificationRules()
		}
	}()
}

func processNotificationRules() {
	var rules []models.NotificationRule
	if err := config.DB.Where("is_active = ?", true).Find(&rules).Error; err != nil {
		log.Printf("Error fetching active notification rules: %v", err)
		return
	}

	for _, rule := range rules {
		if shouldTriggerNotification(rule) {
			createNotification(rule)
			sendEmail(rule)
		}
	}
}

func shouldTriggerNotification(rule models.NotificationRule) bool {
	var requestCount int64

	query := `SELECT COUNT(*) FROM request_logs WHERE hostname = ? AND created_at >= ?`
	err := config.DB.Raw(query, rule.HostName, time.Now().Add(-time.Duration(rule.TimeWindow)*time.Minute)).Scan(&requestCount).Error
	if err != nil {
		log.Printf("Error fetching request count for rule %s: %v", rule.Name, err)
		return false
	}

	return requestCount >= int64(rule.Threshold)
}

func createNotification(rule models.NotificationRule) {
	notification := models.Notification{
		// RuleID:    rule.ID,
		Message: fmt.Sprintf("Threshold exceeded for rule: %s", rule.Name),
		// CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification for rule %s: %v", rule.Name, err)
	}
}

func sendEmail(rule models.NotificationRule) {
	var users []models.NotificationConfig
	if err := config.DB.Where("user_id IN ?", rule.UsersID).Find(&users).Error; err != nil {
		log.Printf("Error fetching users for rule %s: %v", rule.Name, err)
		return
	}

	for _, user := range users {
		fmt.Printf("Email sent!\nSubject: Alert - %s\nTo: %s\nBody: Your WAF detected a threshold breach for rule: %s\n\n", rule.Name, user.Email, rule.Name)
	}
}
