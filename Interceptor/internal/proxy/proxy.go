package proxy

import (
	"fmt"
	"interceptor/internal/logger"
	"interceptor/internal/waf"
	"io"
	"log"
	"net/http"
	// "interceptor/internal/parser"
	"interceptor/internal/error"
	"runtime"
	// "interceptor/internal/ml"
)

const remoteLogServer = "192.168.1.2:514" //get the ip address and port form the server

const targetRedirectIP = "http://127.0.0.1:5000" //get the ip address and port form the server

const proxyPort = ":8080" //get the port form the server

var requestQueue = make(chan *http.Request, 100)

func proxyRequest(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusOK)
		return
	}

	blocked_by_Rule, RuleID, RuleMessage, Action, Status := waf.EvaluateRules(r)
	// request_header_metadata := parser.ParseRequest(r)
	// request_body_metadata:=parser.parseRequestBody(r)
	// blocked_by_ML, logMessage := ml.EvaluateML(request_header_metadata)

	if blocked_by_Rule {
		error.Send403Response(w, RuleID, RuleMessage, Action, Status)
		logger.LogRequest(r, Action, RuleMessage)
		return
	}

	logger.LogRequest(r, "allow", "")

	client := &http.Client{}
	targetURL := fmt.Sprintf("%s%s", targetRedirectIP, r.URL.Path)

	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}
}

func worker() {
	for r := range requestQueue {
		fmt.Printf("Worker processing request: %s %s", r.Method, r.URL)
	}
}

func Starter() {
	if err := waf.InitializeRuleEngine(); err != nil {
		log.Fatalf("Failed to initialize WAF: %v", err)
	}

	workerPoolSize := runtime.NumCPU()
	fmt.Printf("Detected %d logical CPUs. Setting worker pool size to %d.\n", workerPoolSize, workerPoolSize)

	for i := 0; i < workerPoolSize; i++ {
		go worker()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case requestQueue <- r:
		default:
			http.Error(w, "Server too busy. Try again later.", http.StatusServiceUnavailable)
		}
		proxyRequest(w, r)
	})
	err := logger.InitializeLogger(remoteLogServer)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.CloseLogger()

	fmt.Printf("Starting server on port %s", proxyPort)
	if err := http.ListenAndServe(proxyPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
