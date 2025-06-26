package controllers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNotificationRule(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	rules, err := services.FetchNotificationRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification_rules": rules})
}

// Redundant with GetNotificationRule; kept for compatibility
func GetNotificationRules(c *gin.Context) {
	GetNotificationRule(c)
}

func UpdateNotificationRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	var input struct {
		Threshold  int   `json:"threshold" binding:"required"`
		TimeWindow int   `json:"time_window" binding:"required"` // You can enable this when logic is ready
		IsActive   *bool `json:"is_active" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, status := services.UpdateNotificationRule(ruleID, input.Threshold, input.TimeWindow, *input.IsActive)
	c.JSON(status, gin.H{"message": message})
}
