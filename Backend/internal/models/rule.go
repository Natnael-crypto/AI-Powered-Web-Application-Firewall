package models

import "time"

type Rule struct {
	RuleID         string    `json:"rule_id"`
	RuleType       string    `json:"rule_type" binding:"required"`
	RuleMethod     string    `json:"rule_method" binding:"required"`
	RuleDefinition string    `json:"rule_definition" binding:"required"`
	Action         string    `json:"action" binding:"required"`
	ApplicationID  string    `json:"application_id" binding:"required"`
	RuleString     string    `json:"rule_string" binding:"required"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active"`
	Category       string    `json:"category" binding:"required"`
}

// RuleData struct to hold the input JSON data
type RuleData struct {
	RuleID         string `json:"rule_id"`
	RuleType       string `json:"rule_type"`
	RuleMethod     string `json:"rule_method"`
	RuleDefinition string `json:"rule_definition"`
	Action         string `json:"action"`
	Category       string `json:"category"`
}
