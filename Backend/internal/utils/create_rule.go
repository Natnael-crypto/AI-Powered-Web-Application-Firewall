package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

func AddRuleBySystem(input models.RuleInput) error {

	ruleID := generateRuleID()

	input.RuleID = ruleID

	ruleString, err := GenerateRule(input)
	if err != nil {
		fmt.Print("failed to generate rule")
		return err
	}
	rule := models.Rule{
		RuleID:         ruleID,
		RuleDefinition: marshalConditions(input.Conditions),
		Action:         input.Action,
		RuleMethod:     "chained",
		RuleType:       "multiple",
		RuleString:     ruleString,
		CreatedBy:      "System",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsActive:       input.IsActive,
		Category:       input.Category,
	}

	if err := config.DB.Create(&rule).Error; err != nil {
		fmt.Print("failed to create rule")
		return err
	}

	for _, id := range input.ApplicationIDs {
		var app models.Application
		if err := config.DB.Where("application_id = ?", id).First(&app).Error; err != nil {
			fmt.Print("application not found")
			return err
		}
		var rule_to_app models.RuleToApp
		rule_to_app.RuleID = ruleID
		rule_to_app.ApplicationID = id
		if err := config.DB.Create(&rule_to_app).Error; err != nil {
			fmt.Print("failed to create rule to app")
			return err
		}
	}
	config.Change = true
	return nil
}
