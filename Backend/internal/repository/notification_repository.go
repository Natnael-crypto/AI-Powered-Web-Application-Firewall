package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func GetUserByID(userID string) (models.User, error) {
	var user models.User
	err := config.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func CreateNotification(notification models.Notification) error {
	return config.DB.Create(&notification).Error
}

func GetNotificationsByUserID(userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := config.DB.Where("user_id = ?", userID).Find(&notifications).Error
	if err != nil {
		return notifications, err
	}
	return notifications, nil
}

func GetNotificationByID(notificationID string) (models.Notification, error) {
	var notification models.Notification
	err := config.DB.Where("notification_id = ?", notificationID).First(&notification).Error
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func SaveNotification(notification models.Notification) error {
	return config.DB.Save(&notification).Error
}

func DeleteNotification(notificationID string) error {
	return config.DB.Where("notification_id = ?", notificationID).Delete(&models.Notification{}).Error
}
