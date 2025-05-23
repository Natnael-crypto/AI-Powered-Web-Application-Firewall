// controllers/interceptor_controller.go
package controllers

import (
	"backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	if config.InterceptorRunning {
		config.InterceptorRunning = false
		c.JSON(http.StatusOK, gin.H{"message": "Interceptor will start soon."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Interceptor is already running."})
}

func StopInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	if !config.InterceptorRunning {
		config.InterceptorRunning = true
		c.JSON(http.StatusOK, gin.H{"message": "Interceptor will stop soon."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Interceptor is already stopped."})
}

func RestartInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	// Set change to true regardless of the current state
	config.Change = true
	c.JSON(http.StatusOK, gin.H{"message": "Interceptor will restart soon."})
}

func InterceptorCheckState(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"running": config.InterceptorRunning,
		"change":  config.Change,
	})
	// After interceptor fetches the status, reset the change to false
	config.Change = false
}

func MlCheckState(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"model_setting_updated": config.ModelSettingUpdated,
		"select_model":          config.SelectModel,
	})

	config.ModelSettingUpdated = false
	config.SelectModel = false
}
