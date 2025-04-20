package waf

import (
	"fmt"
	"interceptor/internal/utils"
	"io"
	"log"
	"net/http"
	"strconv"

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
func (w *WAF) EvaluateRules(r *http.Request) (bool, int, string, string, int, string) {
	tx := w.engine.NewTransaction() // Correct usage
	defer tx.Close()

	for name, values := range r.Header {
		for _, value := range values {
			if _, err := strconv.ParseFloat(value, 64); err == nil {
				tx.AddRequestHeader(name, value)
				continue
			}

			decodedVal := utils.RecursiveDecode(value, 3)
			tx.AddRequestHeader(name, decodedVal)
		}
	}

	tx.ProcessRequestHeaders()
	url := utils.RecursiveDecode(r.RequestURI, 3)
	tx.ProcessURI(url, r.Method, r.Proto)
	string_body := ""
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		string_body = string(body)
		if err != nil {
			log.Println("error while reading the request body", err)
		}
		fmt.Println(string(body))
		tx.AddPostRequestArgument("body", string(body))
		tx.ProcessRequestBody()
	}

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
				fmt.Println(rule.Message())
			}
		}
	}

	if interruption != nil {
		return true, interruption.RuleID, ruleMessage, interruption.Action, interruption.Status, string_body
	}

	return false, 0, "", "", 0, string_body
}
