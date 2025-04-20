package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddNotificationRule creates a new notification rule
func AddNotificationRule(c *gin.Context) {
	currentUserID := c.GetString("user_id")

	var input struct {
		Name       string   `json:"name" binding:"required"`
		HostName   string   `json:"hostname" binding:"required"`
		ThreatType string   `json:"threat_type" binding:"required"`
		Threshold  int      `json:"threshold" binding:"required"`
		TimeWindow int      `json:"time_window" binding:"required"`
		IsActive   bool     `json:"is_active" binding:"required"`
		UsersID    []string `json:"users_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the hostname exists in the applications table
	var app models.Application
	if err := config.DB.Where("host_name = ?", input.HostName).First(&app).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hostname does not exist in the application list"})
		return
	}

	// Check for duplicate notification rule with same host_name and threat_type
	var existingRule models.NotificationRule
	if err := config.DB.Where("host_name = ? AND threat_type = ?", input.HostName, input.ThreatType).First(&existingRule).Error; err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification rule already exists for this hostname and threat type"})
		return
	}

	// Validate threshold and time window
	if input.Threshold <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold must be greater than 0"})
		return
	}
	if input.TimeWindow <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time window must be greater than 0"})
		return
	}

	// Verify all users exist
	for _, userID := range input.UsersID {
		var user models.User
		if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID: " + userID})
			return
		}
	}

	// Encode users ID to JSON
	jsonUsers, err := json.Marshal(input.UsersID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode users_id"})
		return
	}

	// Create the new notification rule
	rule := models.NotificationRule{
		ID:         uuid.New().String(),
		Name:       input.Name,
		HostName:   input.HostName,
		CreatedBy:  currentUserID,
		ThreatType: input.ThreatType,
		Threshold:  input.Threshold,
		TimeWindow: input.TimeWindow,
		IsActive:   input.IsActive,
		UsersID:    jsonUsers,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification rule added successfully"})
}

// GetNotificationRules retrieves all notification rules
func GetNotificationRules(c *gin.Context) {
	var rules []models.NotificationRule

	if err := config.DB.Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

// UpdateNotificationRule updates an existing notification rule
func UpdateNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")
	currentUserID := c.GetString("user_id")

	// Get the existing rule to check ownership
	var rule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	// Check if current user is the creator of the rule
	if currentUserID != rule.CreatedBy {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var input struct {
		Name       string   `json:"name" binding:"required"`
		HostName   string   `json:"hostname" binding:"required"`
		ThreatType string   `json:"threat_type" binding:"required"`
		Threshold  int      `json:"threshold" binding:"required"`
		TimeWindow int      `json:"time_window" binding:"required"`
		IsActive   bool     `json:"is_active"`
		UsersID    []string `json:"users_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonUsers, err := json.Marshal(input.UsersID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode users_id"})
		return
	}

	// Update fields
	if input.Name != "" {
		rule.Name = input.Name
	}
	if input.HostName != "" {
		rule.HostName = input.HostName
	}
	if input.ThreatType != "" {
		rule.ThreatType = input.ThreatType
	}
	if input.Threshold != 0 {
		rule.Threshold = input.Threshold
	}
	if input.TimeWindow != 0 {
		rule.TimeWindow = input.TimeWindow
	}
	if input.IsActive != false {
		rule.IsActive = input.IsActive
	}
	if input.UsersID != nil {
		rule.UsersID = jsonUsers
	}
	rule.UpdatedAt = time.Now()

	if err := config.DB.Save(&rule).Error; err != nil {
		log.Printf("Error updating notification rule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification rule updated successfully"})
}

// DeleteNotificationRule removes a notification rule
func DeleteNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	currentUserID := c.GetString("user_id")

	// Get the existing rule to check ownership
	var existingRule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&existingRule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	// Check if current user is the creator of the rule
	if currentUserID != existingRule.CreatedBy {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := config.DB.Delete(&existingRule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification rule deleted successfully"})
}

// ToggleNotificationRuleStatus activates or deactivates a notification rule
func ToggleNotificationRuleStatus(c *gin.Context) {
	ruleID := c.Param("rule_id")

	currentUserID := c.GetString("user_id")

	// Get the existing rule to check ownership
	var existingRule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&existingRule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	// Check if current user is the creator of the rule
	if currentUserID != existingRule.CreatedBy {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

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

// AddNotificationConfig creates a new notification config for a user
func AddNotificationConfig(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
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

// GetNotificationConfig retrieves a user's notification config by UserID
func GetNotificationConfig(c *gin.Context) {
	userID := c.Param("user_id")

	var configEntry models.NotificationConfig
	if err := config.DB.Where("user_id = ?", userID).First(&configEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification config not found"})
		return
	}

	currentUserID := c.GetString("user_id")
	// Check if current user is the creator of the rule
	if currentUserID != configEntry.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
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

// UpdateNotificationConfig updates a user's notification config
func UpdateNotificationConfig(c *gin.Context) {
	userID := c.Param("user_id")
	currentUserID := c.GetString("user_id")

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

	// Check if current user is the creator of the rule
	if currentUserID != configEntry.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	configEntry.Email = input.Email

	if err := config.DB.Save(&configEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification config updated successfully"})
}

// DeleteNotificationConfig removes a user's notification config
func DeleteNotificationConfig(c *gin.Context) {
	UserID := c.Param("user_id")
	currentUserID := c.GetString("user_id")

	// Check if current user is the creator of the rule
	if currentUserID != UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if err := config.DB.Where("user_id = ?", UserID).Delete(&models.NotificationConfig{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification config deleted successfully"})
}
