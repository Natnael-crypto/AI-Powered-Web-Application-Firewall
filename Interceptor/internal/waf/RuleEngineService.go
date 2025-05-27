package waf

import (
	"fmt"
	"interceptor/internal/utils"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/corazawaf/coraza/v3"
)

type WAF struct {
	engine coraza.WAF
}

func InitializeRuleEngine(customRule string) (*WAF, error) {
	cfg := coraza.NewWAFConfig().
		WithDirectivesFromFile("./internal/config/crs-setup.conf").
		WithDirectivesFromFile("./internal/config/rules/*.conf").
		WithDirectivesFromFile("./internal/config/custom/" + customRule)

	engine, err := coraza.NewWAF(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize WAF: %v", err)
	}

	return &WAF{engine: engine}, nil
}

func (w *WAF) EvaluateRules(r *http.Request) (bool, int, string, string, int, string) {
	tx := w.engine.NewTransaction()
	defer tx.Close()

	ignoredMessages := []string{
		"Enabling body inspection",
		"Invalid HTTP Request Line",
		"Request Missing a Host Header",
		"",
	}

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

	ruleMessage := ""
	if interruption != nil {
		var matchedMessages []string

		for _, rule := range tx.MatchedRules() {
			if !slices.Contains(ignoredMessages, rule.Message()) {
				matchedMessages = append(matchedMessages, rule.Message())
				if rule.Rule().ID() == interruption.RuleID && rule.Rule().ID() >= 1000000000000000000 {
					ruleMessage = rule.Message()
					break
				}
			}
		}

		if ruleMessage == "" {
			ruleMessage = strings.Join(matchedMessages, " & ")
		}
		return true, interruption.RuleID, ruleMessage, interruption.Action, interruption.Status, string_body
	}
	return false, 0, "", "", 0, string_body
}