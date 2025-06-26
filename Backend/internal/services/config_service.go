package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConfigService fetches the configuration for the system
func GetConfigService() (gin.H, int) {
	var conf models.Conf
	if err := repository.GetConfig(&conf); err != nil {
		return gin.H{"error": "configuration not found"}, http.StatusNotFound
	}
	return gin.H{"data": conf}, http.StatusOK
}

// GetConfigAdminService fetches the configuration for super admins
func GetConfigAdminService() (gin.H, int) {
	var conf models.Conf
	if err := repository.GetConfig(&conf); err != nil {
		return gin.H{"error": "configuration not found"}, http.StatusNotFound
	}
	return gin.H{"data": conf}, http.StatusOK
}

// GetAppConfigService fetches application-specific configuration based on app ID
func GetAppConfigService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	var appConf models.AppConf
	if err := repository.GetAppConfig(applicationID, &appConf); err != nil {
		return gin.H{"error": "App configuration not found"}, http.StatusNotFound
	}
	return gin.H{"data": appConf}, http.StatusOK
}

// CreateConfigService creates a new system-wide configuration
func CreateConfigService(c *gin.Context) (gin.H, int) {
	var input struct {
		ListeningPort   string `json:"listening_port" binding:"numeric"`
		RemoteLogServer string `json:"remote_logServer" binding:"required,max=40"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	newConf := models.Conf{
		ID:              utils.GenerateUUID(),
		ListeningPort:   input.ListeningPort,
		RemoteLogServer: input.RemoteLogServer,
	}

	if err := repository.CreateConfig(newConf); err != nil {
		return gin.H{"error": "Failed to create configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Configuration created successfully", "config": newConf}, http.StatusCreated
}

// UpdateListeningPortService updates the listening port in the configuration
func UpdateListeningPortService(c *gin.Context) (gin.H, int) {
	var input struct {
		ListeningPort string `json:"listening_port"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		return gin.H{"error": "Invalid input format"}, http.StatusBadRequest
	}

	var conf models.Conf
	if err := repository.GetConfig(&conf); err != nil {
		return gin.H{"error": "configuration not found"}, http.StatusNotFound
	}

	conf.ListeningPort = input.ListeningPort
	if err := repository.UpdateConfig(conf); err != nil {
		return gin.H{"error": "failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "configuration updated successfully", "data": conf}, http.StatusOK
}

// UpdateRateLimitService updates the rate limit and related config settings for an app
func UpdateRateLimitService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	var input struct {
		RateLimit  int `json:"rate_limit" binding:"required,min=1"`
		WindowSize int `json:"window_size" binding:"required,min=1"`
		BlockTime  int `json:"block_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		return gin.H{"error": "Invalid input format or missing fields"}, http.StatusBadRequest
	}

	var conf models.AppConf
	if err := repository.GetAppConfig(applicationID, &conf); err != nil {
		return gin.H{"error": "App configuration not found"}, http.StatusNotFound
	}

	conf.RateLimit = input.RateLimit
	conf.WindowSize = input.WindowSize
	conf.BlockTime = input.BlockTime

	if err := repository.UpdateAppConfig(applicationID, conf); err != nil {
		return gin.H{"error": "Failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Configuration updated successfully", "data": conf}, http.StatusOK
}

// UpdateTlsService updates the TLS setting for an app configuration
func UpdateTlsService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	var input struct {
		Tls bool `json:"tls" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		return gin.H{"error": "Invalid input format or missing fields"}, http.StatusBadRequest
	}

	var conf models.AppConf
	if err := repository.GetAppConfig(applicationID, &conf); err != nil {
		return gin.H{"error": "App configuration not found"}, http.StatusNotFound
	}

	conf.Tls = input.Tls
	if err := repository.UpdateAppConfig(applicationID, conf); err != nil {
		return gin.H{"error": "Failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Configuration updated successfully", "data": conf}, http.StatusOK
}

// UpdateDetectBotService updates the detect bot setting for an app configuration
func UpdateDetectBotService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	var input struct {
		DetectBot bool `json:"detect_bot" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		return gin.H{"error": "Invalid input format or missing fields"}, http.StatusBadRequest
	}

	var conf models.AppConf
	if err := repository.GetAppConfig(applicationID, &conf); err != nil {
		return gin.H{"error": "App configuration not found"}, http.StatusNotFound
	}

	conf.DetectBot = input.DetectBot
	if err := repository.UpdateAppConfig(applicationID, conf); err != nil {
		return gin.H{"error": "Failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Configuration updated successfully", "data": conf}, http.StatusOK
}

// UpdateRemoteLogServerService updates the remote log server for the configuration
func UpdateRemoteLogServerService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		RemoteLogServer string `json:"remote_logServer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		return gin.H{"error": "Invalid input format"}, http.StatusBadRequest
	}

	var conf models.Conf
	if err := repository.GetConfig(&conf); err != nil {
		return gin.H{"error": "configuration not found"}, http.StatusNotFound
	}

	conf.RemoteLogServer = input.RemoteLogServer
	if err := repository.UpdateConfig(conf); err != nil {
		return gin.H{"error": "failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "configuration updated successfully", "data": conf}, http.StatusOK
}

func UpdateMaxPostDataSizeService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")

	// Check if user has access to the application
	if c.GetString("role") != "super_admin" {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	// Bind the request data
	var input struct {
		MaxPostDataSize float64 `json:"max_post_data_size"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for MaxPostDataSize update: %v", err)
		return gin.H{"error": "Invalid input format or missing fields"}, http.StatusBadRequest
	}

	// Fetch and update the app configuration
	conf, err := repository.GetAppConfigByID(applicationID)
	if err != nil {
		return gin.H{"error": "App configuration not found"}, http.StatusNotFound
	}

	// Update the MaxPostDataSize field
	conf.MaxPostDataSize = input.MaxPostDataSize

	// Save the updated configuration by dereferencing the conf pointer
	if err := repository.UpdateAppConfig(applicationID, *conf); err != nil {
		return gin.H{"error": "Failed to update configuration"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Configuration updated successfully", "data": conf}, http.StatusOK
}
