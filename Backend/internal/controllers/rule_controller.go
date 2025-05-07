package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
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
	var input models.RuleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appIds := utils.GetAssignedApplicationIDs(c)

	for _, id := range input.ApplicationIDs {
		if !slices.Contains(appIds, id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var app models.Application
	ruleID := generateRuleID()

	for _, id := range input.ApplicationIDs {
		if err := config.DB.Where("application_id = ?", id).First(&app).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		var rule_to_app models.RuleToApp
		rule_to_app.RuleID = ruleID
		rule_to_app.ApplicationID = id
		if err := config.DB.Create(&rule_to_app).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule to app"})
			return
		}
	}

	input.RuleID = ruleID

	ruleString, err := utils.GenerateRule(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate rule"})
		return
	}
	userID := c.GetString("user_id")
	rule := models.Rule{
		RuleID:         ruleID,
		RuleDefinition: marshalConditions(input.Conditions),
		Action:         input.Action,
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

	var rule_to_app []models.RuleToApp
	if err := config.DB.Where("app_id = ?", applicationID).Find(&rule_to_app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rules not found"})
		return
	}

	ruleIDs := make([]string, 0, len(rule_to_app))
	for _, mapping := range rule_to_app {
		ruleIDs = append(ruleIDs, mapping.RuleID)
	}

	var rules []models.Rule
	if err := config.DB.Where("rule_id In ?", ruleIDs).Find(&rules).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rules not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func GetAllRulesAdmin(c *gin.Context) {

	appIDs := utils.GetAssignedApplicationIDs(c)

	var rules []models.Rule
	if err := config.DB.Where("application_id IN ?", appIDs).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func GetOneRule(c *gin.Context) {
	ruleID := c.Param("rule_id")
	user_id := c.GetString("user_id")

	type RuleDefinition struct {
		RuleType       string `json:"rule_type"`
		RuleMethod     string `json:"rule_method"`
		RuleDefinition string `json:"rule_definition"`
	}

	var rule models.Rule
	if err := config.DB.Where("rule_id = ? And created_by = ?", ruleID, user_id).First(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rule"})
		return
	}

	var parsedDefs []RuleDefinition
	if err := json.Unmarshal([]byte(rule.RuleDefinition), &parsedDefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rule definition", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rule":            rule,
		"rule_definition": parsedDefs,
	})
}

func UpdateRule(c *gin.Context) {
	ruleID := c.Param("rule_id")
	user_id := c.GetString("user_id")

	var input models.RuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rule models.Rule
	if err := config.DB.Where("rule_id = ? And created_by = ?", ruleID, user_id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	appIds := utils.GetAssignedApplicationIDs(c)

	for _, id := range input.ApplicationIDs {
		if !slices.Contains(appIds, id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		} else {
			var existing_rule_to_app models.RuleToApp
			if err := config.DB.Where("rule_id = ? And app_id = ?", rule.RuleID, id).First(&existing_rule_to_app).Error; err != nil {
				rule_to_app := models.RuleToApp{
					RuleID:        rule.RuleID,
					ApplicationID: id,
				}
				if err := config.DB.Create(&rule_to_app).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule to app"})
					return
				}
			}
		}
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

	if err := config.DB.Save(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func DeactivateRule(c *gin.Context) {

	ruleID := c.Param("rule_id")
	user_id := c.GetString("user_id")

	var rule models.Rule
	if err := config.DB.Where("rule_id = ? And created_by = ?", ruleID, user_id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	rule.IsActive = false

	if err := config.DB.Save(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func ActivateRule(c *gin.Context) {

	ruleID := c.Param("rule_id")
	user_id := c.GetString("user_id")

	var rule models.Rule
	if err := config.DB.Where("rule_id = ? And created_by = ?", ruleID, user_id).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	rule.IsActive = true

	if err := config.DB.Save(rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
}

func DeleteRule(c *gin.Context) {

	ruleID := c.Param("rule_id")
	user_id := c.GetString("user_id")

	if err := config.DB.Where("rule_id = ? And created_by = ?", ruleID, user_id).Delete(&models.Rule{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule not found"})
		return
	}

	config.Change = true

	var rule_to_app models.RuleToApp
	if err := config.DB.Where("rule_id = ?", ruleID, user_id).Delete(&rule_to_app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rule to application"})
		return
	}

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
