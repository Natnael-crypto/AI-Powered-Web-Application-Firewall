package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetNotificationRule(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var rules []models.NotificationRule

	if err := config.DB.Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

func GetNotificationRules(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var rules []models.NotificationRule

	if err := config.DB.Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

func UpdateNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	var rule models.NotificationRule
	if err := config.DB.Where("id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification rule not found"})
		return
	}

	var input struct {
		Threshold  int    `json:"threshold" binding:"required"`
		TimeWindow int    `json:"time_window" binding:"required"`
		IsActive   bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
