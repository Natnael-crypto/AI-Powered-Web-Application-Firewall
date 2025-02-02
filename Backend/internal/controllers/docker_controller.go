package controllers

import (
	"fmt"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// StartContainer starts the "interceptor" container using docker run command.
func StartContainer(c *gin.Context) {
	// Check if the container exists and is running
	if containerExists() {
		c.JSON(200, gin.H{
			"message": "Interceptor container is already running.",
		})
		return
	}

	// Start the container
	cmd := exec.Command("docker", "run", "-d", "--name", "interceptor", "--network", "waf_network",
		"-e", "BACKENDHOST=backend", "-e", "BACKENDPORT=8080",
		"-e", "MLHOST=cnnapi", "-e", "MLPORT=5000",
		"-p", "80:80", "natnaelcrypto/interceptor:latest")

	err := cmd.Run()
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to start container: %s", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Interceptor container started.",
	})
}

// StopContainer stops the "interceptor" container if it exists.
func StopContainer(c *gin.Context) {
	// Check if the container exists
	if !containerExists() {
		c.JSON(404, gin.H{
			"message": "Interceptor container does not exist.",
		})
		return
	}

	// Stop the container
	cmd := exec.Command("docker", "stop", "interceptor")

	err := cmd.Run()
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to stop container: %s", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Interceptor container stopped.",
	})
}

// RestartContainer restarts the "interceptor" container, or starts it if it doesn't exist.
func RestartContainer(c *gin.Context) {
	// Check if the container exists
	if !containerExists() {
		c.JSON(404, gin.H{
			"message": "Interceptor container does not exist. Starting a new one...",
		})
		StartContainer(c) // Directly calling StartContainer without returning an error
		return
	}

	// Restart the container
	cmd := exec.Command("docker", "restart", "interceptor")

	err := cmd.Run()
	if err != nil {
		c.JSON(500, gin.H{
			"error": fmt.Sprintf("Failed to restart container: %s", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Interceptor container restarted.",
	})
}

// containerExists checks if the "interceptor" container exists and is running.
func containerExists() bool {
	cmd := exec.Command("docker", "ps", "-a", "--filter", "name=interceptor", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil || string(output) == "" {
		return false
	}
	return true
}
