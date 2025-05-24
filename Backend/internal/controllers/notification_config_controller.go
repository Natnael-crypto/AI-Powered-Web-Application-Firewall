package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"log"
	"strings"
	// "backend/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToggleNotificationRuleStatus(c *gin.Context) {
	ruleID := c.Param("rule_id")

	var existingRule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&existingRule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	// var app models.Application
	// if err := config.DB.Where("hostname = ?", existingRule.HostName).First(&app).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch notification rules"})
	// 	return
	// }

	// if c.GetString("role") == "super_admin" {
	// } else {
	// 	appIds := utils.GetAssignedApplicationIDs(c)
	// 	if !utils.HasAccessToApplication(appIds, app.ApplicationID) {
	// 		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
	// 		return
	// 	}
	// }

	existingRule.IsActive = !existingRule.IsActive
	existingRule.UpdatedAt = time.Now()

	if err := config.DB.Save(&existingRule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule status"})
		return
	}

	status := "activated"
	if !existingRule.IsActive {
		status = "deactivated"
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification rule " + status + " successfully"})
}

func AddNotificationConfig(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
	}

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	configEntry := models.NotificationConfig{
		UserID: input.UserID,
		Email:  input.Email,
	}

	if err := config.DB.Create(&configEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification config"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification config added successfully"})
}

func GetNotificationConfig(c *gin.Context) {
	userID := c.GetString("user_id")

	var configEntry models.NotificationConfig
	if err := config.DB.Where("user_id = ?", userID).First(&configEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification config not found"})
		return
	}

	currentUserID := c.GetString("user_id")
	if currentUserID != configEntry.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_config": configEntry})
}

func GetAllNotificationConfig(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var configEntry []models.NotificationConfig
	if err := config.DB.Find(&configEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification config not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_config": configEntry})
}

func GetNotificationConfig_local(c *gin.Context) (string, error) {
	userID := c.Param("user_id")

	var configEntry models.NotificationConfig
	err := config.DB.Where("user_id = ?", userID).First(&configEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("notification config not found")
		}
		return "", fmt.Errorf("failed to retrieve notification config: %w", err)
	}

	return configEntry.Email, nil
}

func UpdateNotificationConfig(c *gin.Context) {

	userID := c.Param("user_id")

	if c.GetString("role") != "super_admin" || c.GetString("user_id") != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var configEntry models.NotificationConfig
	if err := config.DB.Where("user_id = ?", userID).First(&configEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification config not found"})
		return
	}

	configEntry.Email = input.Email

	if err := config.DB.Save(&configEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification config updated successfully"})
}

func DeleteNotificationConfig(c *gin.Context) {
	userID := c.Param("user_id")

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	if err := config.DB.Where("user_id = ?", userID).Delete(&models.NotificationConfig{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification config deleted successfully"})
}

func SaveNotificationSenderConfig(c *gin.Context) {
	var input models.NotificationSender

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	// Check if the notification sender already exists
	var senderConfig models.NotificationSender
	result := config.DB.First(&senderConfig, "email = ?", input.Email) // Find by email

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		// Handle unexpected error
		log.Println("Database Error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// If the sender does not exist, create it
	if result.RowsAffected == 0 {
		if err := config.DB.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification sender"})
			return
		}
	} else {
		// If it exists, update both the email and app password
		if input.AppPassword == strings.Repeat("•", 16) {
			senderConfig.Email = input.Email
		} else {
			senderConfig.Email = input.Email
			senderConfig.AppPassword = input.AppPassword
		}
		if err := config.DB.Save(&senderConfig).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification sender"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sender saved successfully"})
}

func GetNotificationSenderConfig(c *gin.Context) {
	var senderConfig models.NotificationSender

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	if err := config.DB.First(&senderConfig).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification sender config not found"})
	}

	senderConfig.AppPassword = strings.Repeat("•", 16)

	c.JSON(http.StatusOK, senderConfig)
}
