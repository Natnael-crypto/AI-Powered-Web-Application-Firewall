package controllers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ToggleNotificationRuleStatus(c *gin.Context) {
	ruleID := c.Param("rule_id")
	message, statusCode := services.ToggleNotificationRuleStatus(ruleID)
	c.JSON(statusCode, gin.H{"message": message})
}

func AddNotificationConfig(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, statusCode := services.AddNotificationConfig(input.UserID, input.Email)
	c.JSON(statusCode, gin.H{"message": message})
}

func GetNotificationConfig(c *gin.Context) {
	userID := c.GetString("user_id")
	configEntry, err := services.GetNotificationConfig(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notification_config": configEntry})
}

func GetAllNotificationConfig(c *gin.Context) {
	configEntries, err := services.GetAllNotificationConfig()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notification_config": configEntries})
}

func UpdateNotificationConfig(c *gin.Context) {
	userID := c.Param("user_id")
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, statusCode := services.UpdateNotificationConfig(userID, input.Email)
	c.JSON(statusCode, gin.H{"message": message})
}

func DeleteNotificationConfig(c *gin.Context) {
	userID := c.Param("user_id")
	message, statusCode := services.DeleteNotificationConfig(userID)
	c.JSON(statusCode, gin.H{"message": message})
}

func SaveNotificationSenderConfig(c *gin.Context) {
	var input struct {
		Email       string `json:"email"`
		AppPassword string `json:"app_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message, statusCode := services.SaveNotificationSenderConfig(input.Email, input.AppPassword)
	c.JSON(statusCode, gin.H{"message": message})
}

func GetNotificationSenderConfig(c *gin.Context) {
	senderConfig, err := services.GetNotificationSenderConfig()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, senderConfig)
}
