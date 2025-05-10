package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddNotificationRule(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	var input struct {
		Name           string   `json:"name" binding:"required"`
		ThreatType     string   `json:"threat_type" binding:"required"`
		Threshold      int      `json:"threshold" binding:"required"`
		TimeWindow     int      `json:"time_window" binding:"required"`
		IsActive       bool     `json:"is_active" binding:"required"`
		UsersID        []string `json:"users_id" binding:"required"`
		ApplicationsID []string `json:"applicationids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		for _, appId := range input.ApplicationsID {
			if !utils.HasAccessToApplication(appIds, appId) {
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
				return
			}
		}
	}

	if input.Threshold <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "threshold must be greater than 0"})
		return
	}
	if input.TimeWindow <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "time window must be greater than 0"})
		return
	}

	for _, userID := range input.UsersID {
		var user models.User
		if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID: " + userID})
			return
		}
	}

	for _, appId := range input.ApplicationsID {
		var application models.Application
		if err := config.DB.Where("user_id = ?", appId).First(&application).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID: " + appId})
			return
		}
	}

	rule := models.NotificationRule{
		ID:         uuid.New().String(),
		Name:       input.Name,
		CreatedBy:  currentUserID,
		ThreatType: input.ThreatType,
		Threshold:  input.Threshold,
		TimeWindow: input.TimeWindow,
		IsActive:   input.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := config.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification rule"})
		return
	}

	for _, appId := range input.ApplicationsID {

		notificationRuleToApplication := models.NotificationRuleToApplication{
			NotificationRuleID: rule.ID,
			ApplicationID:      appId,
		}

		if err := config.DB.Create(&notificationRuleToApplication).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification rule to application relationship"})
			return
		}
	}

	for _, userId := range input.UsersID {

		notificationRuleToUser := models.NotificationRuleToUser{
			NotificationRuleID: rule.ID,
			UserID:             userId,
		}

		if err := config.DB.Create(&notificationRuleToUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification rule to user relationship"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "notification rule added successfully"})
}

func GetNotificationRule(c *gin.Context) {

	applicationID := c.Param("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var notificationRuleToApplication []models.NotificationRuleToApplication
	if err := config.DB.Where("application_id = ?", applicationID).First(&notificationRuleToApplication).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	var ruleIDs []string
	for _, record := range notificationRuleToApplication {
		ruleIDs = append(ruleIDs, record.NotificationRuleID)
	}

	var rules []models.NotificationRule

	if err := config.DB.Where("id = ?", ruleIDs).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

func GetNotificationRules(c *gin.Context) {

	appIds := utils.GetAssignedApplicationIDs(c)

	var notificationRuleToApplication []models.NotificationRuleToApplication

	if err := config.DB.Where("application_id In ?", appIds).Find(&notificationRuleToApplication).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	var ruleIDs []string
	for _, record := range notificationRuleToApplication {
		ruleIDs = append(ruleIDs, record.NotificationRuleID)
	}

	var rules []models.NotificationRule

	if err := config.DB.Where("id In ?", ruleIDs).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

func UpdateNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")
	currentUserID := c.GetString("user_id")

	var rule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	if currentUserID != rule.CreatedBy {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Name           string   `json:"name" binding:"required"`
		HostName       string   `json:"hostname" binding:"required"`
		ThreatType     string   `json:"threat_type" binding:"required"`
		Threshold      int      `json:"threshold" binding:"required"`
		TimeWindow     int      `json:"time_window" binding:"required"`
		IsActive       bool     `json:"is_active"`
		UsersID        []string `json:"users_id" binding:"required"`
		ApplicationsID []string `json:"applicationids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		rule.Name = input.Name
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
	if input.IsActive {
		rule.IsActive = input.IsActive
	}
	rule.UpdatedAt = time.Now()

	if err := config.DB.Save(&rule).Error; err != nil {
		log.Printf("Error updating notification rule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification rule updated successfully"})
}

func DeleteNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	currentUserID := c.GetString("user_id")

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
