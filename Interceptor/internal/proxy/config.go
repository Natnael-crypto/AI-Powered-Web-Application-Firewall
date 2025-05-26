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
	BlockTime       int     `json:"block_time"`
	HostName        string  `json:"hostname"`
	ApplicationID   string  `json:"application_id"`
	MaxPostDataSize float64 `json:"max_post_data_size" `
	Tls             bool    `json:"tls"`
}

type Application struct {
	ApplicationID   string    `json:"application_id"`
	ApplicationName string    `json:"application_name"`
	Hostname        string    `json:"hostname"`
	IPAddress       string    `json:"ip_address"`
	Port            string    `json:"port"`
	Status          bool      `json:"status"`
	Tls             bool      `json:"tls"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Config          AppConfig `json:"config"`
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

	backendHost := os.Getenv("BACKENDURL")

	configURL := fmt.Sprintf(backendHost + "/interceptor/config")

	req, err := http.NewRequest("GET", configURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Service", "I")

	client := &http.Client{}
	resp, err := client.Do(req)

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

func fetchApplicationConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDURL")

	applicationsURL := fmt.Sprintf(backendHost + "/interceptor/application/")

	req, err := http.NewRequest("GET", applicationsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("X-Service", "I")

	client := &http.Client{}
	resp, err := client.Do(req)

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
			application_config[app.Hostname] = app.Config // set config directly
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

			// Fetch security headers
			appSecurityHeaderURL := fmt.Sprintf(backendHost+"/interceptor/security-headers/%s", app.ApplicationID)

			req, err := http.NewRequest("GET", appSecurityHeaderURL, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %v", err)
			}

			req.Header.Set("X-Service", "I")

			client := &http.Client{}
			securityResp, err := client.Do(req)

			if err != nil {
				log.Printf("Failed to fetch security headers for %s: %v", app.ApplicationID, err)
				continue
			}
			defer securityResp.Body.Close()

			var securityHeadersResult struct {
				Headers []SecurityHeader `json:"security_headers"`
			}

			if err := json.NewDecoder(securityResp.Body).Decode(&securityHeadersResult); err != nil {
				log.Printf("Failed to decode security headers: %v", err)
				continue
			}
			application_security_headers[app.Hostname] = securityHeadersResult.Headers
		}
	}

	fmt.Println("Loaded applications:", applications)
	return nil
}
