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

type AppTOIp struct {
	IPAddress     string
	ApplicationID string
}

func processNotificationRules() {
	var rules []models.NotificationRule
	if err := config.DB.Where("is_active = ?", true).Find(&rules).Error; err != nil {
		log.Printf("Error fetching active notification rules: %v", err)
		return
	}

	for _, rule := range rules {
		timeWindowStart := time.Now().Add(-time.Duration(rule.TimeWindow) * time.Second).UnixMilli()

		targetedApps, _, _ := shouldTriggerNotification(rule, timeWindowStart)

		for appID, count := range targetedApps {

			var Requests []models.Request

			if err := config.DB.Where("application_id = ? And status = true And timestamp >= ? ", appID, timeWindowStart).Find(&Requests).Error; err != nil {
				log.Printf("Error fetching requests for app %s: %v", appID, err)
			}
			var clientIPs []string

			for _, request := range Requests {
				clientIPs = append(clientIPs, request.ClientIP)
			}

			var UserToApp []models.UserToApplication

			if err := config.DB.Where("application_id = ? ", appID).Find(&UserToApp).Error; err != nil {
				log.Printf("Error fetching user to app: %v", err)
			}

			var Application models.Application

			if err := config.DB.Where("application_id = ? ", appID).Find(&Application).Error; err != nil {
				log.Printf("Error fetching application : %v", err)
			}

			var users []models.User

			if err := config.DB.Where("id IN (?)", UserToApp).Find(&users).Error; err != nil {
				log.Printf("Error fetching users: %v", err)
			}

			dashboardMessage := fmt.Sprintf(
				"🚨 %s rule triggered on %s | IP(s): %s | Count: %d",
				rule.Name,
				Application.HostName,
				clientIPs,
				count,
			)

			for _, user := range users {
				createNotification(rule, dashboardMessage)
				emailMessage := fmt.Sprintf(
					`Hello %s,

Your Web Application Firewall has detected suspicious activity.

🛡️ Rule Triggered: %s
📌 Application: %s
🌐 Source IP(s): %s
🔢 Occurrence Count: %d
🕒 Time: %s

Recommended Action: Please review the related logs and ensure appropriate mitigation steps are in place.

Best regards,
WAF Security Monitoring System`,
					user.Username,
					rule.Name,
					Application.HostName,
					clientIPs,
					count,
					time.Now().Format("2006-01-02 15:04:05 MST"),
				)
				sendEmail(rule, user, emailMessage)
			}

		}
	}
}

func shouldTriggerNotification(rule models.NotificationRule, timeWindowStart int64) (map[string]int64, int64, []string) {
	type Result struct {
		ThreatType    string
		AIThreatType  string
		ApplicationID string
		Count         int64
	}

	var results []Result
	targetedApps := make(map[string]int64)
	var totalCount int64
	var matchingThreats []string

	query := `
		SELECT threat_type , ai_threat_type, application_id , COUNT(*) as count
		FROM requests
		WHERE timestamp >= ? AND Status = 'blocked'
		GROUP BY threat_type, ai_threat_type ,application_id
	`

	err := config.DB.Raw(query, timeWindowStart).Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching request data for rule %s: %v", rule.Name, err)
		return targetedApps, totalCount, matchingThreats
	}

	ruleThreat := strings.ToLower(rule.ThreatType)

	for _, res := range results {
		rule_threat := strings.ToLower(res.ThreatType)
		Ml_threat := strings.ToLower(res.AIThreatType)

		if ruleThreat == "*" || strings.Contains(rule_threat, ruleThreat) {
			totalCount += res.Count

			if res.Count >= int64(rule.Threshold) {
				targetedApps[res.ApplicationID] += res.Count
			}

			if !contains(matchingThreats, rule_threat) {
				matchingThreats = append(matchingThreats, rule_threat)
			}
		} else if strings.Contains(rule_threat, Ml_threat) {
			totalCount += res.Count

			if res.Count >= int64(rule.Threshold) {
				targetedApps[res.ApplicationID] += res.Count
			}

			if !contains(matchingThreats, rule_threat) {
				matchingThreats = append(matchingThreats, rule_threat)
			}
		}
	}

	return targetedApps, totalCount, matchingThreats
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

func sendEmail(rule models.NotificationRule, user models.User, message string) {
	var sender models.NotificationSender

	if err := config.DB.First(&sender).Error; err != nil {
		log.Println("Sender email has not been configured:", err)
		return
	}

	var recipient models.NotificationConfig
	if err := config.DB.Where("user_id = ? ", user.UserID).First(&recipient).Error; err != nil {
		log.Println("Failed to retrieve recipient emails:", err)
		return
	}

	SendEmail(recipient, sender, rule, message)

}
