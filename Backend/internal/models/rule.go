package models

import "time"

type Rule struct {
	RuleID         string    `json:"rule_id"`
	RuleType       string    `json:"rule_type" binding:"required"`
	RuleDefinition string    `json:"rule_definition" binding:"required"`
	Action         string    `json:"action" binding:"required"`
	ApplicationID  string    `json:"application_id" binding:"required"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active"`
	Category       string    `json:"category" binding:"required"`
}
