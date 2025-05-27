package background

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	// "encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"slices"

	"github.com/google/uuid"
)

func StartNotificationWatcher() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
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

		fmt.Print(targetedApps)

		fmt.Print(rule.Name)

		for appID, count := range targetedApps {

			var Requests []models.Request

			if err := config.DB.Where("application_id = ? And status = 'blocked' And timestamp >= ? ", appID, timeWindowStart).Find(&Requests).Error; err != nil {
				log.Printf("Error fetching requests for app %s: %v", appID, err)
			}
			var clientIPs []string

			for _, request := range Requests {
				if !slices.Contains(clientIPs, request.ClientIP) {
					clientIPs = append(clientIPs, request.ClientIP)
				}
			}

			var UserToApp []models.UserToApplication

			if err := config.DB.Where("application_id = ? ", appID).Find(&UserToApp).Error; err != nil {
				log.Printf("Error fetching user to app: %v", err)
			}

			var Application models.Application

			if err := config.DB.Where("application_id = ? ", appID).Find(&Application).Error; err != nil {
				log.Printf("Error fetching application : %v", err)
			}

			var userID []string

			for _, user := range UserToApp {
				userID = append(userID, user.UserID)
			}

			var users []models.User

			if err := config.DB.Where("user_id IN ?", userID).Find(&users).Error; err != nil {
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
				createNotification(user, dashboardMessage)
				emailMessage := utils.ComposeEmailMessage(user.Username, rule.Name, Application.HostName, clientIPs, count, timeWindowStart)
				sendEmail(rule, user, emailMessage)
			}
		}
	}
}

func shouldTriggerNotification(rule models.NotificationRule, timeWindowStart int64) (map[string]int64, int64, []string) {
	type Result struct {
		ApplicationID string
		Count         int64
	}

	var results []Result
	targetedApps := make(map[string]int64)
	var totalCount int64
	var matchingThreats []string

	ruleThreat := strings.ToLower(rule.ThreatType)

	var query string
	var args []interface{}

	if ruleThreat == "" || ruleThreat == "*" {
		query = `
			SELECT application_id, COUNT(*) as count
			FROM requests
			WHERE timestamp >= ? AND status = 'blocked'
			GROUP BY application_id
		`
		args = append(args, timeWindowStart)
	} else {
		query = `
			SELECT application_id, COUNT(*) as count
			FROM requests
			WHERE timestamp >= ?
			  AND status = 'blocked'
			  AND (
			      LOWER(threat_type) LIKE ?
			      OR LOWER(ai_threat_type) LIKE ?
			  )
			GROUP BY application_id
		`
		likePattern := "%" + ruleThreat + "%"
		args = append(args, timeWindowStart, likePattern, likePattern)
	}

	err := config.DB.Raw(query, args...).Scan(&results).Error
	if err != nil {
		log.Printf("Error fetching request data for rule %s: %v", rule.Name, err)
		return targetedApps, totalCount, matchingThreats
	}

	for _, res := range results {
		totalCount += res.Count
		if res.Count >= int64(rule.Threshold) {
			targetedApps[res.ApplicationID] += res.Count
		}
	}

	if ruleThreat != "" && ruleThreat != "*" {
		matchingThreats = append(matchingThreats, ruleThreat)
	}

	return targetedApps, totalCount, matchingThreats
}

func createNotification(users models.User, message string) {

	notification := models.Notification{
		NotificationID: uuid.New().String(),
		UserID:         users.UserID,
		Message:        message,
		Timestamp:      time.Now(),
		Status:         false,
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		log.Printf("Error creating notification  %v", err)
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
