package utils

import (
	"backend/internal/models"
	"fmt"
	"strings"
)

func GenerateRule(ruleData models.RuleInput) (string, error) {
	if len(ruleData.Conditions) == 0 {
		return "", fmt.Errorf("no conditions provided")
	}

	validRuleMethods := map[string]bool{
		"regex": true, "streq": true, "contains": true,
		"ipMatch": true, "rx": true, "beginsWith": true,
		"endsWith": true, "eq": true, "pm": true,
	}

	validRuleTypes := map[string]bool{
		"REQUEST_HEADERS": true, "REQUEST_URI": true, "ARGS": true,
		"ARGS_GET": true, "ARGS_POST": true, "REQUEST_COOKIES": true,
		"REQUEST_BODY": true, "XML": true, "JSON": true,
		"REQUEST_METHOD": true, "REQUEST_PROTOCOL": true, "REMOTE_ADDR": true,
	}

	validActions := map[string]bool{
		"deny": true, "log": true, "nolog": true, "pass": true,
		"drop": true, "redirect": true, "capture": true,
		"t:none": true, "t:lowercase": true, "t:normalizePath": true,
		"t:urlDecode": true, "t:compressWhitespace": true,
		"severity:2": true, "severity:3": true, "status:403": true,
	}

	// Validate action tokens (split by comma and check individually)
	actionTokens := strings.Split(ruleData.Action, ",")
	for _, token := range actionTokens {
		t := strings.TrimSpace(token)
		if !validActions[t] && !strings.HasPrefix(t, "t:") && !strings.HasPrefix(t, "status:") && !strings.HasPrefix(t, "severity:") {
			return "", fmt.Errorf("invalid action '%s'", t)
		}
	}

	var ruleBuilder strings.Builder

	for i, cond := range ruleData.Conditions {
		// Validate rule method
		if !validRuleMethods[cond.RuleMethod] {
			return "", fmt.Errorf("invalid rule method '%s'", cond.RuleMethod)
		}

		// Validate rule type
		if !validRuleTypes[cond.RuleType] {
			return "", fmt.Errorf("invalid rule type '%s'", cond.RuleType)
		}

		// Write rules
		if i == 0 {
			// First condition (main rule)
			ruleBuilder.WriteString(fmt.Sprintf("SecRule %s \"@%s %s\" \"id:%s,phase:2,%s,msg:'%s'",
				cond.RuleType,
				cond.RuleMethod,
				cond.RuleDefinition,
				ruleData.RuleID,
				ruleData.Action,
				ruleData.Category,
			))
		} else {
			// Chained condition
			ruleBuilder.WriteString(fmt.Sprintf("\n    chain\n    SecRule %s \"@%s %s\"",
				cond.RuleType,
				cond.RuleMethod,
				cond.RuleDefinition,
			))
		}
	}

	return ruleBuilder.String(), nil
}
