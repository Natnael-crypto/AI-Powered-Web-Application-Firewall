package main

import (
	"fmt"
	"net/http"
	"runtime"

	"ai-powered-waf/rule-engine/internal/config"
	"ai-powered-waf/rule-engine/internal/logger"
	"ai-powered-waf/rule-engine/internal/utils"
	"ai-powered-waf/rule-engine/internal/proxy"
	"ai-powered-waf/rule-engine/internal/waf"
)

func main() {
	// Parse configuration
	config.ParseConfig()

	// Initialize logger
	logger.InitializeLogger()

	// Initialize WAF
	waf.InitializeWAF()

	// Log configuration details
	logger.Info(fmt.Sprintf("Server will run on port: %d", config.Config.Server.Port))
	logger.Info(fmt.Sprintf("Redirecting requests to: %s", config.Config.Server.RedirectURL))

	// Set up worker pool size
	workerPoolSize := runtime.NumCPU()
	logger.Info(fmt.Sprintf("Detected %d logical CPUs. Setting worker pool size to %d.", workerPoolSize, workerPoolSize))

	// Create request queue
	requestQueue := make(chan *http.Request, 100)

	// Start worker pool for processing requests
	for i := 0; i < workerPoolSize; i++ {
		go func() {
			for r := range requestQueue {
				// Parse and log the request
				requestDetails := utils.ParseRequest(r)
				logger.LogRequest(requestDetails)
			}
		}()
	}

	// Define HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case requestQueue <- r:
			// Process request through proxy
			proxy.ProxyRequest(w, r)
			requestDetails := utils.ParseRequest(r)
			logger.LogRequest(requestDetails)
		default:
			http.Error(w, "Server too busy. Try again later.", http.StatusServiceUnavailable)
			logger.Warn("Server too busy. Dropping request.")
		}
	})

	// Start the HTTP server
	port := config.Config.Server.Port
	logger.Info(fmt.Sprintf("Starting server on port %d...", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		logger.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
