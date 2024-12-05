package waf

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/corazawaf/coraza/v3"
)

var WAFInstance coraza.WAF

// InitializeWAF sets up the WAF with custom rules.
func InitializeWAF() error {
	customRules := `
        SecDebugLog ./logs/debug.log
        SecDebugLogLevel 9
        SecRuleEngine On
        SecRequestBodyAccess On
        SecResponseBodyAccess On
				# SQL Injection Rule (common SQL injection patterns)
        SecRule ARGS "@rx (\bselect\b|\binsert\b|\bupdate\b|\bdelete\b|\bdrop\b|\bunion\b|\b--|\b;|\b'\s*or\s*1=1)" \
            "id:1002,phase:2,deny,status:403,log,msg:'SQL Injection Detected: suspicious SQL pattern in input'"

        # XSS Rule (detect <script> tag in arguments)
        SecRule ARGS "@rx <script>" "id:1001,phase:2,deny,status:403,log,msg:'XSS detected: <script> tag in input'"
    `
	tmpFile, err := os.CreateTemp("", "custom_rules_*.conf")
	if err != nil {
		return fmt.Errorf("failed to create temporary rule file: %v", err)
	}
	defer tmpFile.Close()

	if _, err = tmpFile.WriteString(customRules); err != nil {
		return fmt.Errorf("failed to write custom rules: %v", err)
	}

	cfg := coraza.NewWAFConfig().WithDirectivesFromFile(tmpFile.Name())
	WAFInstance, err = coraza.NewWAF(cfg)
	if err != nil {
		return err
	}

	return os.Remove(tmpFile.Name())
}

// RequestChecker processes requests through the WAF.
func RequestChecker(r *http.Request) (bool, string) {
	// Create a transaction to process the request
	tx := WAFInstance.NewTransaction()
	defer tx.Close()

	// Add headers to transaction
	for name, values := range r.Header {
		for _, value := range values {
			tx.AddRequestHeader(name, value)
		}
	}

	tx.ProcessRequestHeaders();
	tx.ProcessURI(r.RequestURI, r.Method, r.Proto)
	tx.ProcessRequestBody()

	// Check if the transaction is interrupted (blocked by a rule)
	interruption := tx.Interruption()
	// fmt.Printf("Is Rule Engine on: %s\n", tx.IsRuleEngineOff())
	matchedRules := tx.MatchedRules()
	// fmt.Print(len(matchedRules))
	var matchedRulesList []string

	if len(matchedRules) > 0 {
		fmt.Println("Matched Rules:")
		ruleIDPrinted := false

		for _, rule := range matchedRules {
			if !ruleIDPrinted {

				// fmt.Printf("Rule ID: %s\n", rule.TransactionID())
				ruleIDPrinted = true
			}
			matchedRulesList = append(matchedRulesList, fmt.Sprintf("Rule Message: %s", rule.Message()))
			// fmt.Printf("Rule Message: %s\n", rule.Message())
			// fmt.Printf("Rule Disruptive: %d\n", rule.Disruptive())
		}
	}

	if interruption != nil {
		// return true, fmt.Sprintf("Request Blocked: %v", interruption)
		return true, strings.Join(matchedRulesList, "\n")
	}

	return false, ""
}
