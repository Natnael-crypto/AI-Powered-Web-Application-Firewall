package controllers

import (
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	message, statusCode := services.StartInterceptor()
	c.JSON(statusCode, gin.H{"message": message})
}

func StopInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	message, statusCode := services.StopInterceptor()
	c.JSON(statusCode, gin.H{"message": message})
}

func RestartInterceptor(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	message, statusCode := services.RestartInterceptor()
	c.JSON(statusCode, gin.H{"message": message})
}

func InterceptorCheckState(c *gin.Context) {
	state := services.GetInterceptorState()
	c.JSON(http.StatusOK, gin.H{
		"running": state.Running,
		"change":  state.Change,
	})
}

func MlCheckState(c *gin.Context) {
	state := services.GetMlState()
	c.JSON(http.StatusOK, gin.H{
		"model_setting_updated": state.ModelSettingUpdated,
		"select_model":          state.SelectModel,
	})
}
