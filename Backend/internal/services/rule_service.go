package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ServiceError struct {
	Message string
	Status  int
}

func AddRuleService(c *gin.Context) *ServiceError {
	var input models.RuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return &ServiceError{err.Error(), http.StatusBadRequest}
	}

	appIds := utils.GetAssignedApplicationIDs(c)
	for _, id := range input.ApplicationIDs {
		if !contains(appIds, id) {
			return &ServiceError{"insufficient privileges", http.StatusForbidden}
		}
	}

	ruleID := generateRuleID()
	input.RuleID = ruleID

	ruleString, err := utils.GenerateRule(input)
	if err != nil {
		return &ServiceError{"failed to generate rule", http.StatusInternalServerError}
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

	if err := repository.CreateRule(&rule, input.ApplicationIDs); err != nil {
		return &ServiceError{err.Error(), http.StatusInternalServerError}
	}

	config.Change = true
	c.JSON(http.StatusCreated, gin.H{"message": "rule added successfully", "rule": rule})
	return nil
}

func GetRulesService(c *gin.Context) ([]models.Rule, *ServiceError) {
	appID := c.Param("application_id")
	ruleIDs, err := repository.GetRuleIDsByApp(appID)
	if err != nil {
		return nil, &ServiceError{"rules not found", http.StatusNotFound}
	}

	rules, err := repository.GetActiveRulesByIDs(ruleIDs)
	if err != nil {
		return nil, &ServiceError{"rules not found", http.StatusNotFound}
	}
	return rules, nil
}

func GetAllRulesAdminService(c *gin.Context) ([]gin.H, *ServiceError) {
	appIDs := utils.GetAssignedApplicationIDs(c)
	mappings, err := repository.GetRuleToAppsByAppIDs(appIDs)
	if err != nil {
		return nil, &ServiceError{"failed to fetch rule mappings", http.StatusInternalServerError}
	}

	ruleAppIDs := make([]string, 0, len(mappings))
	ruleToAppsMap := make(map[string][]string) // rule_id -> []application_id

	for _, mapping := range mappings {
		ruleAppIDs = append(ruleAppIDs, mapping.RuleID)
		ruleToAppsMap[mapping.RuleID] = append(ruleToAppsMap[mapping.RuleID], mapping.ApplicationID)
	}

	rules, err := repository.GetRulesByIDs(ruleAppIDs)
	if err != nil {
		return nil, &ServiceError{"failed to fetch rules", http.StatusInternalServerError}
	}

	apps, err := repository.GetAppsByIDs(appIDs)
	if err != nil {
		return nil, &ServiceError{"failed to fetch applications", http.StatusInternalServerError}
	}

	appMap := make(map[string]models.Application)
	for _, a := range apps {
		appMap[a.ApplicationID] = a
	}

	var response []gin.H
	for _, rule := range rules {
		var conditions []models.RuleCondition
		if err := json.Unmarshal([]byte(rule.RuleDefinition), &conditions); err != nil {
			return nil, &ServiceError{"Failed to parse rule definition", http.StatusInternalServerError}
		}
		appObjs := []models.ApplicationOptions{}
		for _, aid := range ruleToAppsMap[rule.RuleID] {
			if app, ok := appMap[aid]; ok {
				appObjs = append(appObjs, models.ApplicationOptions{
					HostName:      app.HostName,
					ApplicationID: app.ApplicationID,
				})
			}
		}
		response = append(response, gin.H{
			"rule_id": rule.RuleID, "rule_definition": conditions,
			"rule_type": rule.RuleType, "rule_method": rule.RuleMethod,
			"action": rule.Action, "rule_string": rule.RuleString,
			"created_at": rule.CreatedAt, "updated_at": rule.UpdatedAt,
			"is_active": rule.IsActive, "category": rule.Category,
			"applications": appObjs,
		})
	}
	return response, nil
}

func GetOneRuleService(c *gin.Context) (gin.H, *ServiceError) {
	ruleID := c.Param("rule_id")
	userID := c.GetString("user_id")
	rule, err := repository.GetRuleByIDAndUser(ruleID, userID)
	if err != nil {
		return nil, &ServiceError{"rule not found", http.StatusNotFound}
	}

	var parsedDefs []models.RuleCondition
	if err := json.Unmarshal([]byte(rule.RuleDefinition), &parsedDefs); err != nil {
		return nil, &ServiceError{"Failed to parse rule definition", http.StatusInternalServerError}
	}

	apps, err := repository.GetAppsByRuleID(ruleID)
	if err != nil {
		return nil, &ServiceError{"Failed to fetch applications", http.StatusInternalServerError}
	}

	return gin.H{
		"rule":            rule,
		"rule_definition": parsedDefs,
		"applications":    apps,
	}, nil
}

func UpdateRuleService(c *gin.Context) *ServiceError {
	ruleID := c.Param("rule_id")
	userID := c.GetString("user_id")

	var input models.RuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return &ServiceError{err.Error(), http.StatusBadRequest}
	}

	rule, err := repository.GetRuleByIDAndUser(ruleID, userID)
	if err != nil {
		return &ServiceError{"rule not found", http.StatusNotFound}
	}

	appIds := utils.GetAssignedApplicationIDs(c)
	for _, id := range input.ApplicationIDs {
		if !contains(appIds, id) {
			return &ServiceError{"insufficient privileges", http.StatusForbidden}
		}
	}

	if err := repository.DeleteRuleToApps(ruleID); err != nil {
		return &ServiceError{"failed to clear old mappings", http.StatusInternalServerError}
	}

	if err := repository.AddRuleToApps(ruleID, input.ApplicationIDs); err != nil {
		return &ServiceError{"failed to add new mappings", http.StatusInternalServerError}
	}

	input.RuleID = ruleID
	ruleString, err := utils.GenerateRule(input)
	if err != nil {
		return &ServiceError{"failed to regenerate rule", http.StatusInternalServerError}
	}

	rule.RuleDefinition = marshalConditions(input.Conditions)
	rule.Action = input.Action
	rule.IsActive = input.IsActive
	rule.Category = input.Category
	rule.RuleString = ruleString
	rule.UpdatedAt = time.Now()

	if err := repository.SaveRule(&rule); err != nil {
		return &ServiceError{"failed to update rule", http.StatusInternalServerError}
	}

	config.Change = true
	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
	return nil
}

func ToggleRuleActivation(c *gin.Context, activate bool) *ServiceError {
	ruleID := c.Param("rule_id")
	userID := c.GetString("user_id")

	rule, err := repository.GetRuleByIDAndUser(ruleID, userID)
	if err != nil {
		return &ServiceError{"rule not found", http.StatusNotFound}
	}

	rule.IsActive = activate
	if err := repository.SaveRule(&rule); err != nil {
		return &ServiceError{"failed to update rule", http.StatusInternalServerError}
	}

	config.Change = true
	c.JSON(http.StatusOK, gin.H{"message": "rule updated successfully", "rule": rule})
	return nil
}

func DeleteRuleService(c *gin.Context) *ServiceError {
	ruleID := c.Param("rule_id")
	userID := c.GetString("user_id")

	if err := repository.DeleteRule(ruleID, userID); err != nil {
		return &ServiceError{"rule not found", http.StatusNotFound}
	}

	config.Change = true
	c.JSON(http.StatusOK, gin.H{"message": "rule deleted successfully"})
	return nil
}

func GetRuleMetadataService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"actions": mapsKeys(validActions),
		"methods": mapsKeys(validRuleMethods),
		"types":   mapsKeys(validRuleTypes),
	})
}

func marshalConditions(conds []models.RuleCondition) string {
	data, _ := json.Marshal(conds)
	return string(data)
}

func generateRuleID() string {
	min := int64(1_000_000_000_000_000_000)
	max := int64(9_223_372_036_854_775_807)
	return strconv.FormatInt(rand.Int63n(max-min)+min, 10)
}

func contains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

func mapsKeys[T comparable](m map[T]any) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Static metadata
var (
	validRuleTypes = map[string]any{
		"REQUEST_HEADERS": nil, "REQUEST_URI": nil, "ARGS": nil,
		"ARGS_GET": nil, "ARGS_POST": nil, "REQUEST_COOKIES": nil,
		"REQUEST_BODY": nil, "XML": nil, "JSON": nil,
		"REQUEST_METHOD": nil, "REQUEST_PROTOCOL": nil, "REMOTE_ADDR": nil,
	}
	validRuleMethods = map[string]any{
		"regex": nil, "streq": nil, "contains": nil,
		"ipMatch": nil, "rx": nil, "beginsWith": nil,
		"endsWith": nil, "eq": nil, "pm": nil,
	}
	validActions = map[string]any{
		"deny": nil, "drop": nil, "pass": nil, "log": nil,
		"redirect": nil, "proxy": nil, "auditlog": nil,
		"status": nil, "tag": nil, "msg": nil,
		"capture": nil, "setvar": nil,
	}
)
