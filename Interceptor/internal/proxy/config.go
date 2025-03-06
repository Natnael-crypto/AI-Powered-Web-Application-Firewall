package proxy

import (
	"encoding/json"
	"fmt"
	"interceptor/internal/waf"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	ID              string `json:"ID"`
	ListeningPort   string `json:"ListeningPort"`
	RemoteLogServer string `json:"RemoteLogServer"`
	RateLimit       int    `json: "rate_limit"`
	WindowSize      int    `json: "window_size"`
}

type Application struct {
	ApplicationID   string `json:"application_id"`
	ApplicationName string `json:"application_name"`
	Hostname        string `json:"hostname"`
	IPAddress       string `json:"ip_address"`
	Port            string `json:"port"`
	Status          bool   `json:"status"`
	Tls             bool   `json:tls`
}

type Cert struct {
	HostName string `json:"hostname"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json: key_path`
}

var CertApp struct {
	Certs []Cert     `json:"certs"`
	mu    sync.Mutex // Ensures thread safety when modifying Certs
}

var (
	remoteLogServer string
	proxyPort       string
	applications    map[string]string   // Maps hostname -> "IP:Port"
	wafInstances    map[string]*waf.WAF // Maps hostname -> WAF instance
	Apps            map[string]Application
	configLock      sync.RWMutex
	appsLock        sync.RWMutex
	rateLimit       int
	windowSize      int
)

// fetchConfig retrieves the configuration from the remote API
func fetchConfig() error {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDHOST")
	if backendHost == "" {
		return fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		return fmt.Errorf("BACKENDPORT environment variable is not set")
	}

	configURL := fmt.Sprintf("http://%s:%s/config", backendHost, backendPort)

	resp, err := http.Get(configURL)
	if err != nil {
		return fmt.Errorf("failed to fetch config: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Config Config `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode config: %v", err)
	}

	configLock.Lock()
	remoteLogServer = result.Config.RemoteLogServer
	proxyPort = ":" + result.Config.ListeningPort
	rateLimit = result.Config.RateLimit
	windowSize = result.Config.WindowSize
	configLock.Unlock()

	fmt.Printf("Loaded config: LogServer=%s, ProxyPort=%s\n", remoteLogServer, proxyPort)
	return nil
}

// fetchApplications retrieves the list of applications and updates the hostname mapping
func fetchApplications() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDHOST")
	if backendHost == "" {
		return fmt.Errorf("BACKENDHOST environment variable is not set")
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		return fmt.Errorf("BACKENDPORT environment variable is not set")
	}

	applicationsURL := fmt.Sprintf("http://%s:%s/application/", backendHost, backendPort)

	resp, err := http.Get(applicationsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch applications: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Applications []Application `json:"applications"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode applications: %v", err)
	}

	appsLock.Lock()
	applications = make(map[string]string)
	wafInstances = make(map[string]*waf.WAF)
	Apps = make(map[string]Application)
	appsLock.Unlock()

	for _, app := range result.Applications {
		if app.Status {
			address := app.IPAddress + ":" + app.Port

			appsLock.Lock()
			applications[app.Hostname] = address
			Apps[app.Hostname] = app
			appsLock.Unlock()

			if app.Tls {
				certPath, keyPath, err := fetchCert(app.ApplicationID)
				if err != nil {
					log.Printf("Failed to get cert: %v", err)
					continue
				}

				CertApp.mu.Lock()
				CertApp.Certs = append(CertApp.Certs, Cert{
					HostName: app.Hostname,
					CertPath: certPath,
					KeyPath:  keyPath,
				})
				CertApp.mu.Unlock()
			}

			rulesResponse, err := FetchRules(app.ApplicationID)
			if err != nil {
				log.Printf("Error fetching rules: %v", err)
				continue
			}

			fileName, err := WriteRuleToFile(app.ApplicationID, rulesResponse.Rules)
			if err != nil {
				log.Printf("Error writing rules to file: %v", err)
				continue
			}

			wafInstance, err := waf.InitializeRuleEngine(fileName)
			if err != nil {
				log.Printf("Error initializing WAF for application %s: %v", app.Hostname, err)
				continue
			}

			appsLock.Lock()
			wafInstances[app.Hostname] = wafInstance
			appsLock.Unlock()
		}
	}

	fmt.Println("Loaded applications:", applications)
	return nil
}
