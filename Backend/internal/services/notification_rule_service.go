package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"log"
	"net/http"
	"time"
)

func FetchNotificationRules() ([]models.NotificationRule, error) {
	return repository.GetAllNotificationRules()
}

func UpdateNotificationRule(ruleID string, threshold int, timeWindow int, isActive bool) (string, int) {
	rule, err := repository.GetNotificationRuleByID(ruleID)
	if err != nil {
		return "notification rule not found", http.StatusNotFound
	}

	if threshold != 0 {
		rule.Threshold = threshold
	}

	// Enable this logic if you want to use timeWindow later
	// if timeWindow != 0 {
	// 	rule.TimeWindow = timeWindow
	// }

	rule.IsActive = isActive
	rule.UpdatedAt = time.Now()

	if err := repository.SaveNotificationRule(rule); err != nil {
		log.Printf("Error updating notification rule: %v", err)
		return "failed to update notification rule", http.StatusInternalServerError
	}

	return "notification rule updated successfully", http.StatusOK
}
