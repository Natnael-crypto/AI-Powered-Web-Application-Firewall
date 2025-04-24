package proxy

import (
	"encoding/json"
	"fmt"
	"interceptor/internal/utils"
	"interceptor/internal/waf"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ID              string `json:"ID"`
	ListeningPort   string `json:"listening_port"`
	RemoteLogServer string `json:"remote_logServer"`
}

type AppConfig struct {
	ID              string  `json:"id"`
	RateLimit       int     `json:"rate_limit"`
	DetectBot       bool    `json:"detect_bot"`
	WindowSize      int     `json:"window_size"`
	HostName        string  `json:"host_name"`
	ApplicationID   string  `json:"application_id"`
	MaxPostDataSize float64 `json:"max_post_data_size" `
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

type SecurityHeader struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	ApplicationID string    `json:"application_id" gorm:"not null"`
	HeaderName    string    `json:"header_name" gorm:"unique;not null"`
	HeaderValue   string    `json:"header_value" gorm:"not null"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Cert struct {
	HostName string `json:"hostname"`
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

var CertApp struct {
	Certs []Cert `json:"certs"`
	mu    sync.Mutex
}

var (
	remoteLogServer              string
	proxyPort                    string
	applications                 map[string]string
	wafInstances                 map[string]*waf.WAF
	Apps                         map[string]Application
	application_config           map[string]AppConfig
	application_security_headers map[string][]SecurityHeader
	configLock                   sync.RWMutex
	appsLock                     sync.RWMutex
)

var WsKey string

func fetchConfig() error {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	WsKey = os.Getenv("WSKEY")

	if WsKey == "" {
		return fmt.Errorf("WsKey environment variable is not set")
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
	configLock.Unlock()

	fmt.Printf("Loaded config: LogServer=%s, ProxyPort=%s\n", remoteLogServer, proxyPort)

	return nil
}

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
	application_config = make(map[string]AppConfig)
	application_security_headers = make(map[string][]SecurityHeader)

	appsLock.Unlock()

	for _, app := range result.Applications {
		if app.Status {
			address := app.IPAddress + ":" + app.Port

			appsLock.Lock()
			applications[app.Hostname] = address
			Apps[app.Hostname] = app
			appsLock.Unlock()

			if app.Tls {
				certPath, keyPath, err := utils.FetchCert(app.ApplicationID)
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

			rulesResponse, err := utils.FetchRules(app.ApplicationID)
			if err != nil {
				log.Printf("Error fetching rules: %v", err)
				continue
			}

			fileName, err := utils.WriteRuleToFile(app.ApplicationID, rulesResponse.Rules)
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

	for _, app := range Apps {
		appConfigURL := fmt.Sprintf("http://%s:%s/config/get-app-config/%s", backendHost, backendPort, app.ApplicationID)
		appSecurityHeaderURL := fmt.Sprintf("http://%s:%s/security-headers/%s", backendHost, backendPort, app.ApplicationID)
		resp, err = http.Get(appConfigURL)

		var appConfigResult struct {
			AppConfig AppConfig `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&appConfigResult); err != nil {
			return fmt.Errorf("failed to decode config: %v", err)
		}

		application_config[app.Hostname] = appConfigResult.AppConfig

		securityResp, err := http.Get(appSecurityHeaderURL)
		if err != nil {
			return fmt.Errorf("failed to fetch security headers for %s: %v", app.ApplicationID, err)
		}
		defer securityResp.Body.Close()

		var securityHeadersResult struct {
			Headers []SecurityHeader `json:"security_headers"`
		}

		if err := json.NewDecoder(securityResp.Body).Decode(&securityHeadersResult); err != nil {
			return fmt.Errorf("failed to decode security headers: %v", err)
		}
		application_security_headers[app.Hostname] = securityHeadersResult.Headers

	}

	fmt.Println("Loaded applications:", applications)
	return nil
}
