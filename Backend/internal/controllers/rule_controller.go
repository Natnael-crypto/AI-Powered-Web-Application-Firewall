package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddRule adds a new rule to the application by a superadmin or assigned admin
func AddRule(c *gin.Context) {
	var input struct {
		RuleType       string `json:"rule_type" binding:"required"`
		RuleDefinition string `json:"rule_definition" binding:"required"`
		Action         string `json:"action" binding:"required"`
		ApplicationID  string `json:"application_id" binding:"required"`
		IsActive       bool   `json:"is_active"`
		Category       string `json:"category" binding:"required"`
	}

	// Parse the input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the application exists
	var app models.Application
	if err := config.DB.Where("application_id = ?", input.ApplicationID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Check if the user has permission to add rules for the application
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	if userRole != "super_admin" {
		var userToApp models.UserToApplication
		if err := config.DB.Where("user_id = ? AND application_id = ?", userID, input.ApplicationID).First(&userToApp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	// Create the rule
	rule := models.Rule{
		RuleID:         uuid.New().String(),
		RuleType:       input.RuleType,
		RuleDefinition: input.RuleDefinition,
		Action:         input.Action,
		ApplicationID:  input.ApplicationID,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsActive:       input.IsActive,
		Category:       input.Category,
	}

	// Save to the database
	if err := config.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "rule added successfully", "rule": rule})
}

// GetRules fetches all rules for a given application
func GetRules(c *gin.Context) {
	applicationID := c.Param("application_id")

	var rules []models.Rule
	if err := config.DB.Where("application_id = ?", applicationID).Find(&rules).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rules not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// UpdateRule updates an existing rule
func UpdateRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	var input struct {
		RuleType       string `json:"rule_type" binding:"required"`
		RuleDefinition string `json:"rule_definition" binding:"required"`
		Action         string `json:"action" binding:"required"`
		IsActive       bool   `json:"is_active"`
		Category       string `json:"category" binding:"required"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the rule exists
	var rule models.Rule
	if err := config.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	// Update the rule details
	rule.RuleType = input.RuleType
	rule.RuleDefinition = input.RuleDefinition
	rule.Action = input.Action
	rule.IsActive = input.IsActive
	rule.Category = input.Category
	rule.UpdatedAt = time.Now()

	// Save the updated rule, specifying the `rule_id`
	if err := config.DB.Model(&models.Rule{}).Where("rule_id = ?", ruleID).Updates(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

// DeleteRule deletes a rule by its ID
func DeleteRule(c *gin.Context) {
	ruleID := c.Param("rule_id")

	if err := config.DB.Where("rule_id = ?", ruleID).Delete(&models.Rule{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted successfully"})
}
