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

// CreateNotification creates a new notification
func CreateNotification(input models.NotificationInput) (string, error) {

	// Input validation already handled by binding tags
	if input.UserID == "" || input.NotificationType == "" || input.Message == "" || input.Severity == "" {
		return "", errors.New("missing required fields")
	}

	// Validate severity level
	validSeverities := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validSeverities[input.Severity] {
		return "", errors.New("invalid severity level")
	}

	// Validate notification type
	validTypes := map[string]bool{
		"alert":   true,
		"warning": true,
		"info":    true,
	}
	if !validTypes[input.NotificationType] {
		return "", errors.New("invalid notification type")
	}

	// Check if the user exists
	var user models.User
	if err := config.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		return "", errors.New("user not found")

	}

	// Create the notification
	notification := models.Notification{
		NotificationID:   uuid.New().String(),
		UserID:           input.UserID,
		Message:          input.Message,
		Timestamp:        time.Now(),
		Status:           input.Status,
	}

	// Save to the database
	if err := config.DB.Create(&notification).Error; err != nil {
		return "", errors.New("failed to create notification")
	}

	return "notification created successfully", nil
}

// GetNotifications fetches all notifications for a given user
func GetNotifications(c *gin.Context) {
	userId := c.Param("user_id")

	// Get the current user's role and ID from the context
	currentUserID := c.GetString("user_id") // Assuming the user ID is set in the context

	// Check if the current user is the same as the requested user or if the user is a super admin
	if currentUserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Fetch the notifications for the user
	var notifications []models.Notification
	if err := config.DB.Where("user_id = ?", userId).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notifications not found"})
		return
	}

	// Return the notifications as a JSON response
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// UpdateNotification updates an existing notification
func UpdateNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")

	// Get the notification to check ownership
	var existingNotification models.Notification
	if err := config.DB.Where("notification_id = ?", notificationID).First(&existingNotification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	// Get current user's ID from context
	currentUserID := c.GetString("user_id")

	// Check if the current user owns this notification or is a super admin
	if currentUserID != existingNotification.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var input struct {
		Status bool `json:"status" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the notification details
	existingNotification.Status = input.Status

	// Save the updated notification, specifying the `notification_id`
	if err := config.DB.Model(&models.Notification{}).Where("notification_id = ?", notificationID).Updates(existingNotification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "notification updated successfully", "notification": existingNotification})
}

// DeleteNotification deletes a notification by its ID
func DeleteNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")

	var existingNotification models.Notification
	if err := config.DB.Where("notification_id = ?", notificationID).First(&existingNotification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	// Get current user's ID from context
	currentUserID := c.GetString("user_id")

	// Check if the current user owns this notification or is a super admin
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
