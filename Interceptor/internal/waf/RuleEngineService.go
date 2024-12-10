package waf

import (
	"fmt"
	"net/http"

	"github.com/corazawaf/coraza/v3"
)

var waf coraza.WAF
var err error

func InitializeRuleEngine() error {

	cfg := coraza.NewWAFConfig().WithDirectivesFromFile("./internal/config/coreruleset/crs-setup.conf").WithDirectivesFromFile("./internal/config/coreruleset/rules/*.conf")
	waf, err = coraza.NewWAF(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize WAF: %v", err)
	}
	return nil
}

func EvaluateRules(r *http.Request) (bool, int, string, string, int) {
	tx := waf.NewTransaction()
	defer tx.Close()

	for name, values := range r.Header {
		for _, value := range values {
			tx.AddRequestHeader(name, value)
		}
	}

	tx.ProcessRequestHeaders()
	tx.ProcessURI(r.RequestURI, r.Method, r.Proto)
	tx.ProcessRequestBody()
	interruption := tx.Interruption()

	ignoredMessages := map[string]bool{
		"Enabling body inspection":      true,
		"Invalid HTTP Request Line":     true,
		"Request Missing a Host Header": true,
	}

	matchedRules := tx.MatchedRules()
	totalRules := len(matchedRules)
	rule_message := ""
	if totalRules > 1 {
		fmt.Println("Matched Rules:")
		ruleIDPrinted := false

		for i, rule := range matchedRules {

			if i == totalRules-1 {
				continue
			}

			if ignoredMessages[rule.Message()] {
				continue
			}

			if !ruleIDPrinted {
				fmt.Printf("Rule ID: %s\n", rule.TransactionID())
				ruleIDPrinted = true
			}

			if len(rule.Message()) > 0 {
				rule_message = rule.Message()
			}
		}
	}

	if interruption != nil {
		return true, interruption.RuleID, rule_message, interruption.Action, interruption.Status
	}

	return false, 0, "", "", 0
}
