package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func GetAllNotificationRules() ([]models.NotificationRule, error) {
	var rules []models.NotificationRule
	err := config.DB.Find(&rules).Error
	return rules, err
}

// func GetNotificationRuleByID(ruleID string) (models.NotificationRule, error) {
// 	var rule models.NotificationRule
// 	err := config.DB.Where("id = ?", ruleID).First(&rule).Error
// 	return rule, err
// }

func SaveNotificationRule(rule models.NotificationRule) error {
	return config.DB.Save(&rule).Error
}
