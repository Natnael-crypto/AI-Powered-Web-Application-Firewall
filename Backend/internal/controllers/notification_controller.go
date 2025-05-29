package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	var user models.User
	if err := config.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		return "", errors.New("user not found")

	}

	notification := models.Notification{
		NotificationID: uuid.New().String(),
		UserID:         input.UserID,
		Message:        input.Message,
		Timestamp:      time.Now(),
		Status:         input.Status,
	}

	if err := config.DB.Create(&notification).Error; err != nil {
		return "", errors.New("failed to create notification")
	}

	return "notification created successfully", nil
}

func GetNotifications(c *gin.Context) {
	currentUserID := c.GetString("user_id")

	var notifications []models.Notification
	if err := config.DB.Where("user_id = ?", currentUserID).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notifications not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func UpdateNotification(c *gin.Context) {

	currentUserID := c.GetString("user_id")
	notificationID := c.Param("notification_id")

	var existingNotification models.Notification
	if err := config.DB.Where("notification_id = ?", notificationID).First(&existingNotification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	if currentUserID != existingNotification.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	existingNotification.Status = !existingNotification.Status

	if err := config.DB.Save(existingNotification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification updated successfully"})
}

func DeleteNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")

	var existingNotification models.Notification
	if err := config.DB.Where("notification_id = ?", notificationID).First(&existingNotification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	currentUserID := c.GetString("user_id")

	if currentUserID != existingNotification.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := config.DB.Where("notification_id = ?", notificationID).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification deleted successfully"})
}
