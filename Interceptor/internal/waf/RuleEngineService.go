package waf

import (
	"fmt"
	"net/http"

	"github.com/corazawaf/coraza/v3"
)

// WAF structure to hold Coraza WAF instance
type WAF struct {
	engine coraza.WAF // Directly store coraza.WAF instead of a pointer
}

// InitializeRuleEngine initializes a new WAF instance with custom rules
func InitializeRuleEngine(customRule string) (*WAF, error) {
	cfg := coraza.NewWAFConfig().
		WithDirectivesFromFile("./internal/config/crs-setup.conf").
		WithDirectivesFromFile("./internal/config/rules/*.conf").
		WithDirectivesFromFile("./internal/config/custom/" + customRule)

	engine, err := coraza.NewWAF(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WAF: %v", err)
	}

	return &WAF{engine: engine}, nil // Wrap it in a struct
}

// EvaluateRules processes incoming requests and applies WAF rules
func (w *WAF) EvaluateRules(r *http.Request) (bool, int, string, string, int) {
	tx := w.engine.NewTransaction() // Correct usage
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
	ruleMessage := ""

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
				ruleMessage = rule.Message()
			}
		}
	}

	if interruption != nil {
		return true, interruption.RuleID, ruleMessage, interruption.Action, interruption.Status
	}

	return false, 0, "", "", 0
}
