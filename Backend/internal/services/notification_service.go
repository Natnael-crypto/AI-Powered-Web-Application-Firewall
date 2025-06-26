package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

func CreateNotification(input models.NotificationInput) (string, error) {
	if input.UserID == "" || input.NotificationType == "" || input.Message == "" || input.Severity == "" {
		return "", errors.New("missing required fields")
	}

	validSeverities := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validSeverities[input.Severity] {
		return "", errors.New("invalid severity level")
	}

	validTypes := map[string]bool{
		"alert":   true,
		"warning": true,
		"info":    true,
	}
	if !validTypes[input.NotificationType] {
		return "", errors.New("invalid notification type")
	}

	_, err := repository.GetUserByID(input.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	notification := models.Notification{
		NotificationID: uuid.New().String(),
		UserID:         input.UserID,
		Message:        input.Message,
		Timestamp:      time.Now(),
		Status:         input.Status,
	}

	err = repository.CreateNotification(notification)
	if err != nil {
		return "", errors.New("failed to create notification")
	}

	return "notification created successfully", nil
}

func GetNotifications(userID string) ([]models.Notification, error) {
	return repository.GetNotificationsByUserID(userID)
}

func UpdateNotification(notificationID, userID string) (string, error) {
	notification, err := repository.GetNotificationByID(notificationID)
	if err != nil {
		return "", errors.New("notification not found")
	}

	if notification.UserID != userID {
		return "", errors.New("access denied")
	}

	notification.Status = !notification.Status
	err = repository.SaveNotification(notification)
	if err != nil {
		return "", errors.New("failed to update notification")
	}

	return "notification updated successfully", nil
}

func UpdateNotificationBatch(notificationIDs []string, userID string) (string, error) {
	for _, notificationID := range notificationIDs {
		notification, err := repository.GetNotificationByID(notificationID)
		if err != nil {
			return "", errors.New("notification not found")
		}

		if notification.UserID != userID {
			return "", errors.New("access denied")
		}

		notification.Status = true
		err = repository.SaveNotification(notification)
		if err != nil {
			return "", errors.New("failed to update notification")
		}
	}

	return "notifications updated successfully", nil
}

func DeleteNotification(notificationID, userID string) (string, error) {
	notification, err := repository.GetNotificationByID(notificationID)
	if err != nil {
		return "", errors.New("notification not found")
	}

	if notification.UserID != userID {
		return "", errors.New("access denied")
	}

	err = repository.DeleteNotification(notificationID)
	if err != nil {
		return "", errors.New("failed to delete notification")
	}

	return "notification deleted successfully", nil
}
