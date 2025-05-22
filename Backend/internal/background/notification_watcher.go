package background

import (
	"backend/internal/config"
	"backend/internal/models"
	// "encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

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
		triggeredIPs, _, _ := shouldTriggerNotification(rule)

		for ip, count := range triggeredIPs {
			message := fmt.Sprintf("Threshold exceeded for rule: %s\nSource IP: %s\nCount: %d", rule.Name, ip, count)
			createNotification(rule, message)
			sendEmail(rule, message)
		}
	}
}

func shouldTriggerNotification(rule models.NotificationRule) (map[string]int64, int64, []string) {
	type Result struct {
		IPAddress  string
		ThreatType string
		Count      int64
	}

	var results []Result
	triggeredIPs := make(map[string]int64)
	var totalCount int64
	var matchingThreats []string

	timeWindowStart := time.Now().Add(-time.Duration(rule.TimeWindow) * time.Second).UnixMilli()

	query := `
		SELECT client_ip, threat_type, COUNT(*) as count
		FROM requests
		WHERE application_name = ? AND timestamp >= ? AND Status = 'blocked'
		GROUP BY client_ip, threat_type
	`

	err := config.DB.Raw(query, timeWindowStart).Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching request data for rule %s: %v", rule.Name, err)
		return triggeredIPs, totalCount, matchingThreats
	}

	ruleThreat := strings.ToLower(rule.ThreatType)

	for _, res := range results {
		threat := strings.ToLower(res.ThreatType)

		if ruleThreat == "*" || strings.Contains(threat, ruleThreat) {
			totalCount += res.Count

			if res.Count >= int64(rule.Threshold) {
				triggeredIPs[res.IPAddress] += res.Count
			}

			if !contains(matchingThreats, threat) {
				matchingThreats = append(matchingThreats, threat)
			}
		}
	}

	return triggeredIPs, totalCount, matchingThreats
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func createNotification(rule models.NotificationRule, message string) {
	var userIDs []string


	for _, userID := range userIDs {
		notification := models.Notification{
			NotificationID: uuid.New().String(),
			UserID:         userID,
			Message:        message,
			Timestamp:      time.Now(),
			Status:         false,
		}

		if err := config.DB.Create(&notification).Error; err != nil {
			log.Printf("Error creating notification for rule %s for user %s: %v", rule.Name, userID, err)
		}
	}
}

func sendEmail(rule models.NotificationRule, message string) {
	var sender models.NotificationSender

	if err := config.DB.First(&sender).Error; err != nil {
		log.Println("Sender email has not been configured:", err)
		return
	}

	var recipients []models.NotificationConfig
	if err := config.DB.Find(&recipients).Error; err != nil {
		log.Println("Failed to retrieve recipient emails:", err)
		return
	}

	for _, recipient := range recipients {
		SendEmail(recipient, sender, rule, message)
	}
}
