package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func generateRuleID() string {
	rand.Seed(time.Now().UnixNano())
	number := rand.Int63n(1000000000000000000)
	return strconv.FormatInt(number, 10)
}

func marshalConditions(conditions []models.RuleCondition) string {
	b, err := json.Marshal(conditions)
	if err != nil {
		return ""
	}
	return string(b)
}

func AddRule(c *gin.Context) {
	applicationID := c.Query("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var input models.RuleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var app models.Application
	if err := config.DB.Where("application_id = ?", input.ApplicationID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	if userRole != "super_admin" {
		var userToApp models.UserToApplication
		if err := config.DB.Where("user_id = ? AND application_id = ?", userID, input.ApplicationID).First(&userToApp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	ruleID := generateRuleID()
	input.RuleID = ruleID

	ruleString, err := utils.GenerateRule(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate rule"})
		return
	}

	rule := models.Rule{
		RuleID:         ruleID,
		RuleDefinition: marshalConditions(input.Conditions),
		Action:         input.Action,
		ApplicationID:  input.ApplicationID,
		RuleMethod:     "chained",
		RuleType:       "multiple",
		RuleString:     ruleString,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsActive:       input.IsActive,
		Category:       input.Category,
	}

	if err := config.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule"})
		return
	}
	config.Change = true

	c.JSON(http.StatusCreated, gin.H{"message": "rule added successfully", "rule": rule})
}

func GetRules(c *gin.Context) {

	applicationID := c.Param("application_id")

	var rules []models.Rule
	if err := config.DB.Where("application_id = ?", applicationID).Find(&rules).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rules not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func GetRulesAdmin(c *gin.Context) {

	applicationID := c.Param("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var rules []models.Rule
	if err := config.DB.Where("application_id = ?", applicationID).Find(&rules).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rules not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func UpdateRule(c *gin.Context) {
	applicationID := c.Query("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	ruleID := c.Param("rule_id")

	var input models.RuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.Rule
	if err := config.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	input.RuleID = ruleID
	ruleString, err := utils.GenerateRule(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to regenerate rule"})
		return
	}

	rule.RuleDefinition = marshalConditions(input.Conditions)
	rule.Action = input.Action
	rule.IsActive = input.IsActive
	rule.Category = input.Category
	rule.UpdatedAt = time.Now()
	rule.RuleString = ruleString
	rule.RuleMethod = "chained"
	rule.RuleType = "multiple"

	if err := config.DB.Model(&models.Rule{}).Where("rule_id = ?", ruleID).Updates(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func DeactivateRule(c *gin.Context) {
	applicationID := c.Query("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	ruleID := c.Param("rule_id")

	var rule models.Rule
	if err := config.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	rule.IsActive = false

	if err := config.DB.Model(&models.Rule{}).Where("rule_id = ?", ruleID).Updates(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func ActivateRule(c *gin.Context) {

	applicationID := c.Query("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	ruleID := c.Param("rule_id")

	var rule models.Rule
	if err := config.DB.Where("rule_id = ?", ruleID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	rule.IsActive = true

	if err := config.DB.Model(&models.Rule{}).Where("rule_id = ?", ruleID).Updates(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func DeleteRule(c *gin.Context) {

	applicationID := c.Query("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	ruleID := c.Param("rule_id")

	if err := config.DB.Where("rule_id = ?", ruleID).Delete(&models.Rule{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule deleted successfully"})
}

var validRuleTypes = map[string]bool{
	"REQUEST_HEADERS": true, "REQUEST_URI": true, "ARGS": true,
	"ARGS_GET": true, "ARGS_POST": true, "REQUEST_COOKIES": true,
	"REQUEST_BODY": true, "XML": true, "JSON": true,
	"REQUEST_METHOD": true, "REQUEST_PROTOCOL": true, "REMOTE_ADDR": true,
}

var validRuleMethods = map[string]bool{
	"regex": true, "streq": true, "contains": true,
	"ipMatch": true, "rx": true, "beginsWith": true,
	"endsWith": true, "eq": true, "pm": true,
}

var validActions = map[string]bool{
	"deny": true, "drop": true, "pass": true, "log": true,
	"redirect": true, "proxy": true, "auditlog": true,
	"status": true, "tag": true, "msg": true,
	"capture": true, "setvar": true,
}

func GetRuleMetadata(c *gin.Context) {
	actions := make([]string, 0, len(validActions))
	for k := range validActions {
		actions = append(actions, k)
	}

	methods := make([]string, 0, len(validRuleMethods))
	for k := range validRuleMethods {
		methods = append(methods, k)
	}

	types := make([]string, 0, len(validRuleTypes))
	for k := range validRuleTypes {
		types = append(types, k)
	}

	c.JSON(http.StatusOK, gin.H{
		"actions": actions,
		"methods": methods,
		"types":   types,
	})
}
