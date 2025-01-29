package utils

import (
	"backend/internal/models"
	"fmt"
	"strings"
)

// Define the possible rule methods
var validRuleMethods = []string{
	"regex", "streq", "contains", "ipMatch", "rx", "beginsWith", "endsWith", "eq", "pm",
}

// generateRule function to generate WAF rules
func GenerateRule(ruleData models.RuleData) (string, error) {
	// Validate the rule method
	validMethod := false
	for _, method := range validRuleMethods {
		if ruleData.RuleMethod == method {
			validMethod = true
			break
		}
	}

	if !validMethod {
		return "", fmt.Errorf("invalid rule method '%s'. Valid methods are: %s", ruleData.RuleMethod, strings.Join(validRuleMethods, ", "))
	}

	// Generate the rule based on the method
	var rule string
	switch ruleData.RuleMethod {
	case "regex":
		rule = fmt.Sprintf("SecRule %s \"@rx %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "streq":
		rule = fmt.Sprintf("SecRule %s \"@streq %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "contains":
		rule = fmt.Sprintf("SecRule %s \"@contains %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "ipMatch":
		rule = fmt.Sprintf("SecRule %s \"@ipMatch %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "rx":
		rule = fmt.Sprintf("SecRule %s \"@rx %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "beginsWith":
		rule = fmt.Sprintf("SecRule %s \"@beginsWith %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "endsWith":
		rule = fmt.Sprintf("SecRule %s \"@endsWith %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "eq":
		rule = fmt.Sprintf("SecRule %s \"@eq %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	case "pm":
		rule = fmt.Sprintf("SecRule %s \"@pm %s\" \"id:%s,phase:2,%s,status:403,msg:'%s'\"", ruleData.RuleType, ruleData.RuleDefinition, ruleData.RuleID, ruleData.Action, ruleData.Category)
	default:
		return "", fmt.Errorf("invalid rule method '%s'", ruleData.RuleMethod)
	}

	return rule, nil
}
