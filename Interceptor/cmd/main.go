package main

import (
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
	Change bool `json:"change"`
}

func checkForChange() bool {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to environment variables")
	}

	backendHost := os.Getenv("BACKENDHOST")

	backendPort := os.Getenv("BACKENDPORT")

	changeURL := "http://" + backendHost + ":" + backendPort + "/change"

	resp, err := http.Get(changeURL)
	if err != nil {
		log.Printf("Error checking for change: %v", err)
		return false
	}
	defer resp.Body.Close()

	var changeResp ChangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&changeResp); err != nil {
		log.Printf("Failed to decode change response: %v", err)
		return false
	}
	return changeResp.Change
}

func main() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			changed := checkForChange()
			if changed {
				log.Println("Change detected. Restarting proxy...")

				// Restart the proxy by spawning new process
				cmd := exec.Command("go", "run", "./cmd")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err := cmd.Start()
				if err != nil {
					log.Fatalf("Failed to restart proxy: %v", err)
				}

				log.Println("New proxy process started with PID:", cmd.Process.Pid)
				os.Exit(0) // kill this one
			}
		}
	}()

	// Continue normal proxy startup
	proxy.Starter()
}
