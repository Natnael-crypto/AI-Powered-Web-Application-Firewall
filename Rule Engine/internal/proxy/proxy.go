package proxy

import (
	"fmt"
	"io" 
	"net/http"
	"ai-powered-waf/rule-engine/internal/config"
	"ai-powered-waf/rule-engine/internal/logger"
	"ai-powered-waf/rule-engine/internal/waf"
)

// ProxyRequest forwards requests to the target server and returns the response.
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Check request using WAF
	blocked, logMessage := waf.RequestChecker(r)
	if blocked {
		logger.Warn(fmt.Sprintf("Request blocked by WAF: %s", logMessage))
		http.Error(w, logMessage, http.StatusForbidden)
		return
	}

	// Construct the target URL
	targetURL := fmt.Sprintf("%s%s", config.Config.Server.RedirectURL, r.URL.Path)
	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}

	// Forward the request to the target server
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create proxy request: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Copy headers to the proxy request
	copyHeaders(r.Header, proxyReq.Header)

	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to reach target server: %v", err))
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers and status code from response
	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)

	// Copy the response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		logger.Error(fmt.Sprintf("Failed to copy response body: %v", err))
	}
	logger.Info(fmt.Sprintf("Successfully proxied request to: %s", targetURL))
}

// copyHeaders copies headers from source to destination.
func copyHeaders(src, dst http.Header) {
	for name, values := range src {
		for _, value := range values {
			dst.Add(name, value)
		}
	}
}
