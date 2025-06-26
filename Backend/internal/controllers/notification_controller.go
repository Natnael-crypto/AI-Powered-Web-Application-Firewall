package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNotification(c *gin.Context) {
	var input models.NotificationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := services.CreateNotification(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func GetNotifications(c *gin.Context) {
	currentUserID := c.GetString("user_id")

	notifications, err := services.GetNotifications(currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func UpdateNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")
	currentUserID := c.GetString("user_id")

	message, err := services.UpdateNotification(notificationID, currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func UpdateNotificationBatch(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	var input struct {
		IDS []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := services.UpdateNotificationBatch(input.IDS, currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DeleteNotification(c *gin.Context) {
	notificationID := c.Param("notification_id")
	currentUserID := c.GetString("user_id")

	message, err := services.DeleteNotification(notificationID, currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
