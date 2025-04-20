package models

import "time"

type Rule struct {
	RuleID         string    `json:"rule_id" gorm:"primaryKey"`
	RuleType       string    `json:"rule_type" gorm:"not null" `
	RuleMethod     string    `json:"rule_method" gorm:"not null" `
	RuleDefinition string    `json:"rule_definition" gorm:"not null"`
	Action         string    `json:"action" gorm:"unique;not null" `
	ApplicationID  string    `json:"application_id" gorm:"not null" `
	RuleString     string    `json:"rule_string" gorm:"not null" `
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active"`
	Category       string    `json:"category" gorm:"not null" `
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

type RuleCondition struct {
	RuleType       string `json:"rule_type" `
	RuleMethod     string `json:"rule_method" `
	RuleDefinition string `json:"rule_definition" `
}

type RuleInput struct {
	Conditions    []RuleCondition `json:"conditions" ` // list of conditions for chaining
	Action        string          `json:"action" `
	ApplicationID string          `json:"application_id" `
	IsActive      bool            `json:"is_active" `
	Category      string          `json:"category" `
	RuleID        string          `json:"rule_id"`
}
