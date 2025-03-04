package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartInterceptor starts the "interceptor" container using docker run command.
func StartInterceptor(c *gin.Context) {
	resp, err := http.Get("http://localhost:80/interceptor/startinterceptor")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	c.Status(resp.StatusCode)
}

// StopInterceptor stops the "interceptor" container if it exists.
func StopInterceptor(c *gin.Context) {
	resp, err := http.Get("http://localhost:80/interceptor/stopinterceptor")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	c.Status(resp.StatusCode)
}

// RestartInterceptor restarts the "interceptor" container, or starts it if it doesn't exist.
func RestartInterceptor(c *gin.Context) {
	resp, err := http.Get("http://localhost:80/interceptor/restartinterceptor")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	c.Status(resp.StatusCode)
}
