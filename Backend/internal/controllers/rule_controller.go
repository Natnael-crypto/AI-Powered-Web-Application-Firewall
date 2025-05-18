package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
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

	ruleID := generateRuleID()

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

	for _, id := range input.ApplicationIDs {
		var app models.Application
		if err := config.DB.Where("application_id = ?", id).First(&app).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
			return
		}
		fmt.Print("Passed")
		var rule_to_app models.RuleToApp
		rule_to_app.RuleID = ruleID
		rule_to_app.ApplicationID = id
		if err := config.DB.Create(&rule_to_app).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule to app"})
			return
		}
	}
	config.Change = true
	c.JSON(http.StatusCreated, gin.H{"message": "rule added successfully", "rule": rule})
}

func GetRules(c *gin.Context) {

	applicationID := c.Param("application_id")

	var rule_to_app []models.RuleToApp
	if err := config.DB.Where("application_id = ?", applicationID).Find(&rule_to_app).Error; err != nil {
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

	var ruleApp []models.RuleToApp
	if err := config.DB.Where("application_id IN ?", appIDs).Find(&ruleApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules to apps mapping"})
		return
	}

	// Collect unique RuleIDs
	ruleAppIDs := make([]string, 0, len(ruleApp))
	ruleToAppsMap := make(map[string][]string) // rule_id -> []application_id

	for _, mapping := range ruleApp {
		ruleAppIDs = append(ruleAppIDs, mapping.RuleID)
		ruleToAppsMap[mapping.RuleID] = append(ruleToAppsMap[mapping.RuleID], mapping.ApplicationID)
	}

	// Fetch rules by rule_id
	var rules []models.Rule
	if err := config.DB.Where("rule_id IN ?", ruleAppIDs).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}

	// Fetch applications
	var applications []models.Application
	if err := config.DB.Where("application_id IN ?", appIDs).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	// Map application_id -> application object
	appMap := make(map[string]models.Application)
	for _, app := range applications {
		appMap[app.ApplicationID] = app
	}

	// Build final response
	var response []gin.H

	for _, rule := range rules {
		var conditions []models.RuleCondition
		if err := json.Unmarshal([]byte(rule.RuleDefinition), &conditions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse rule definition"})
			return
		}

		appIDs := ruleToAppsMap[rule.RuleID]
		appObjs := make([]models.Application, 0, len(appIDs))
		for _, appID := range appIDs {
			if app, ok := appMap[appID]; ok {
				appObjs = append(appObjs, app)
			}
		}

		response = append(response, gin.H{
			"rule_id":        rule.RuleID,
			"rule_type":      rule.RuleType,
			"rule_method":    rule.RuleMethod,
			"rule_definition": conditions,
			"action":         rule.Action,
			"rule_string":    rule.RuleString,
			"created_at":     rule.CreatedAt,
			"updated_at":     rule.UpdatedAt,
			"is_active":      rule.IsActive,
			"category":       rule.Category,
			"applications":   appObjs,
		})
	}

	c.JSON(http.StatusOK, gin.H{"rules": response})
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

	var ruleApp []models.RuleToApp
	if err := config.DB.Where("rule_id = ?", ruleID).Find(&ruleApp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rules"})
		return
	}

	ruleAppIDs := make([]string, 0, len(ruleApp))
	for _, mapping := range ruleApp {
		ruleAppIDs = append(ruleAppIDs, mapping.ApplicationID)
	}

	var applications []models.Application

	if err := config.DB.Where("application_id IN ?", ruleAppIDs).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rule":            rule,
		"rule_definition": parsedDefs,
		"applications":    applications,
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

	if err := config.DB.Where("rule_id = ?", rule.RuleID).Delete(&models.RuleToApp{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rule to app mappings"})
		return
	}

	for _, id := range input.ApplicationIDs {
		if !slices.Contains(appIds, id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		} else {
			var count int64
			if err := config.DB.Model(&models.RuleToApp{}).
				Where("rule_id = ? AND application_id = ?", rule.RuleID, id).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check rule to app mapping"})
				return
			}

			if count < 1 {
				ruleToApp := models.RuleToApp{
					RuleID:        rule.RuleID,
					ApplicationID: id,
				}
				if err := config.DB.Create(&ruleToApp).Error; err != nil {
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
