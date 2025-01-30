package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateNotification creates a new notification
func CreateNotification(c *gin.Context) {
	var input struct {
		UserID           string `json:"user_id" binding:"required"`
		NotificationType string `json:"notification_type" binding:"required"`
		Message          string `json:"message" binding:"required"`
		Status           bool   `json:"status" binding:"required"`
		Severity         string `json:"severity" binding:"required"`
	}

	// Parse the input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.User
	if err := config.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Create the notification
	notification := models.Notification{
		NotificationID:   uuid.New().String(),
		UserID:           input.UserID,
		NotificationType: input.NotificationType,
		Message:          input.Message,
		Timestamp:        time.Now(),
		Status:           input.Status,
		Severity:         input.Severity,
	}

	// Save to the database
	if err := config.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification created successfully", "notification": notification})
}

// GetNotifications fetches all notifications for a given user
func GetNotifications(c *gin.Context) {
	userId := c.Param("user_id")

	// Get the current user's role and ID from the context
	currentUserID := c.GetString("user_id") // Assuming the user ID is set in the context
	currentUserRole := c.GetString("role")  // Assuming the user role is set in the context

	// Check if the current user is the same as the requested user or if the user is a super admin
	if currentUserID != userId && currentUserRole != "super_admin" {
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

	var input struct {
		Status bool `json:"status" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the notification exists
	var notification models.Notification
	if err := config.DB.Where("notification_id = ?", notificationID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	// Update the notification details
	notification.Status = input.Status

	// Save the updated notification, specifying the `notification_id`
	if err := config.DB.Model(&models.Notification{}).Where("notification_id = ?", notificationID).Updates(notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "notification updated successfully", "notification": notification})
}

// DeleteNotification deletes a notification by its ID
func DeleteNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")

	if err := config.DB.Where("notification_id = ?", notificationID).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification deleted successfully"})
}
