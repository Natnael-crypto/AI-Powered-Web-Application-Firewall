package main

import (
	"crypto/tls"
	"encoding/json"
	"interceptor/internal/proxy"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

type ChangeResponse struct {
	Change  bool `json:"change"`
	Running bool `json:"running"`
}

var changed bool

func checkForChange() (bool, bool) {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDURL")

	changeURL := backendHost + "/interceptor/is-running"

	req, err := http.NewRequest("GET", changeURL, nil)
	if err != nil {
		return true, false
	}

	req.Header.Set("X-Service", "I")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error checking for change: %v", err)
		return true, false
	}
	defer resp.Body.Close()

	var changeResp ChangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&changeResp); err != nil {
		log.Printf("Failed to decode change response: %v", err)
		return true, false
	}
	return changeResp.Running, changeResp.Change
}

func main() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			proxy.MaintenanceMode, changed = checkForChange()
			if changed {
				log.Println("Change detected. Restarting proxy...")

				cmd := exec.Command("go", "run", "./cmd")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Start()
				if err != nil {
					log.Fatalf("Failed to restart proxy: %v", err)
				}

				log.Println("New proxy process started with PID:", cmd.Process.Pid)
				os.Exit(0)
			}
		}
	}()

	proxy.Starter()
}
