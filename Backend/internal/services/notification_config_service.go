package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"net/http"
)

func ToggleNotificationRuleStatus(ruleID string) (string, int) {
	existingRule, err := repository.GetNotificationRuleByID(ruleID)
	if err != nil {
		return "notification rule not found", http.StatusNotFound
	}

	existingRule.IsActive = !existingRule.IsActive
	if err := repository.SaveNotificationRule(existingRule); err != nil {
		return "failed to update rule status", http.StatusInternalServerError
	}

	status := "activated"
	if !existingRule.IsActive {
		status = "deactivated"
	}

	return "notification rule " + status + " successfully", http.StatusOK
}

func AddNotificationConfig(userID, email string) (string, int) {
	configEntry := models.NotificationConfig{
		UserID: userID,
		Email:  email,
	}

	if err := repository.CreateNotificationConfig(configEntry); err != nil {
		return "failed to create notification config", http.StatusInternalServerError
	}

	return "notification config added successfully", http.StatusCreated
}

func GetNotificationConfig(userID string) (models.NotificationConfig, error) {
	return repository.GetNotificationConfigByUserID(userID)
}

func GetAllNotificationConfig() ([]models.NotificationConfig, error) {
	return repository.GetAllNotificationConfig()
}

func UpdateNotificationConfig(userID, email string) (string, int) {
	configEntry, err := repository.GetNotificationConfigByUserID(userID)
	if err != nil {
		return "notification config not found", http.StatusNotFound
	}

	configEntry.Email = email
	if err := repository.SaveNotificationConfig(configEntry); err != nil {
		return "failed to update notification config", http.StatusInternalServerError
	}

	return "notification config updated successfully", http.StatusOK
}

func DeleteNotificationConfig(userID string) (string, int) {
	if err := repository.DeleteNotificationConfigByUserID(userID); err != nil {
		return "failed to delete notification config", http.StatusInternalServerError
	}
	return "notification config deleted successfully", http.StatusOK
}

func SaveNotificationSenderConfig(email, appPassword string) (string, int) {
	senderConfig, err := repository.GetNotificationSenderConfig()
	if err != nil {
		// If the sender does not exist, create it
		senderConfig = models.NotificationSender{
			Email:       email,
			AppPassword: appPassword,
		}
		if err := repository.CreateNotificationSenderConfig(senderConfig); err != nil {
			return "Failed to create notification sender", http.StatusInternalServerError
		}
	} else {
		// If exists, update email and app password
		senderConfig.Email = email
		senderConfig.AppPassword = appPassword
		if err := repository.SaveNotificationSenderConfig(senderConfig); err != nil {
			return "Failed to update notification sender", http.StatusInternalServerError
		}
	}

	return "Notification sender saved successfully", http.StatusOK
}

func GetNotificationSenderConfig() (models.NotificationSender, error) {
	return repository.GetNotificationSenderConfig()
}
