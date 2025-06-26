package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func GetNotificationRuleByID(ruleID string) (models.NotificationRule, error) {
	var rule models.NotificationRule
	err := config.DB.Where("id = ?", ruleID).First(&rule).Error
	if err != nil {
		return rule, err
	}
	return rule, nil
}

// func SaveNotificationRule(rule models.NotificationRule) error {
// 	return config.DB.Save(&rule).Error
// }

func CreateNotificationConfig(configEntry models.NotificationConfig) error {
	return config.DB.Create(&configEntry).Error
}

func GetNotificationConfigByUserID(userID string) (models.NotificationConfig, error) {
	var configEntry models.NotificationConfig
	err := config.DB.Where("user_id = ?", userID).First(&configEntry).Error
	if err != nil {
		return configEntry, err
	}
	return configEntry, nil
}

func GetAllNotificationConfig() ([]models.NotificationConfig, error) {
	var configEntries []models.NotificationConfig
	err := config.DB.Find(&configEntries).Error
	if err != nil {
		return configEntries, err
	}
	return configEntries, nil
}

func SaveNotificationConfig(configEntry models.NotificationConfig) error {
	return config.DB.Save(&configEntry).Error
}

func DeleteNotificationConfigByUserID(userID string) error {
	return config.DB.Where("user_id = ?", userID).Delete(&models.NotificationConfig{}).Error
}

func GetNotificationSenderConfig() (models.NotificationSender, error) {
	var senderConfig models.NotificationSender
	err := config.DB.First(&senderConfig).Error
	if err != nil {
		return senderConfig, err
	}
	return senderConfig, nil
}

func CreateNotificationSenderConfig(senderConfig models.NotificationSender) error {
	return config.DB.Create(&senderConfig).Error
}

func SaveNotificationSenderConfig(senderConfig models.NotificationSender) error {
	return config.DB.Save(&senderConfig).Error
}
