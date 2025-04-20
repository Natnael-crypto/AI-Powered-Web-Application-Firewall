package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetConfig retrieves a configuration by ID
func GetConfig(c *gin.Context) {
	var conf models.Conf
	if err := config.DB.First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}
	fmt.Println(conf)
	c.JSON(http.StatusOK, gin.H{"data": conf})
}

func GetAppConfig(c *gin.Context) {
	var appConf models.AppConf
	application_id := c.Param("application_id")

	if err := config.DB.Where("application_id=?", application_id).First(&appConf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App configuration not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": appConf})

}

func CreateAppConfig(c *gin.Context) {

	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ApplicationID   string  `json:application_id`
		RateLimit       int     `json:"rate_limit" binding:"numeric"`
		WindowSize      int     `json:"window_size" binding:"numeric"`
		DetectBot       bool    `json:detect_bot`
		HostName        string  `json:"hostname" binding:"required,max=40"`
		MaxPostDataSize float64 `json:"max_post_data_size" `
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAppConfig models.AppConf

	if err := config.DB.Where("application_id=?", input.ApplicationID).First(&existingAppConfig).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A configuration entry already exists for this Application"})
		return
	}

	var existingApplication models.Application

	if err := config.DB.Where("application_id=?", input.ApplicationID).First(&existingApplication).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Application does not exist"})
		return
	}

	newAppConf := models.AppConf{
		ID:              uuid.New().String(),
		ApplicationID:   input.ApplicationID,
		RateLimit:       input.RateLimit,
		WindowSize:      input.WindowSize,
		DetectBot:       input.DetectBot,
		HostName:        input.HostName,
		MaxPostDataSize: input.MaxPostDataSize,
	}

	if err := config.DB.Create(&newAppConf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	config.Change = true

	c.JSON(http.StatusCreated, gin.H{"message": "Configuration created successfully", "config": newAppConf})
}

func CreateAppConfigLocal(conf models.AppConf) error {

	if err := config.DB.Create(&conf).Error; err != nil {
		return err
	}
	return nil
}

// CreateConfig handles the creation of a new config entry
func CreateConfig(c *gin.Context) {

	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ListeningPort   string `json:"listening_port" binding:"numeric"`
		RemoteLogServer string `json:"remote_logServer" binding:"required,max=40"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a config entry already exists
	var existingConfig models.Conf
	if err := config.DB.First(&existingConfig).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "A configuration entry already exists"})
		return
	}

	// Create a new config entry
	newConf := models.Conf{
		ID:              uuid.New().String(),
		ListeningPort:   input.ListeningPort,
		RemoteLogServer: input.RemoteLogServer,
	}

	if err := config.DB.Create(&newConf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Configuration created successfully", "config": newConf})
}

// UpdateConfig updates an existing configuration
func UpdateListeningPort(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ListeningPort string `json:"listening_port"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging purposes
		log.Printf("Error binding JSON for application update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var conf models.Conf
	if err := config.DB.First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	conf.ListeningPort = input.ListeningPort

	if err := config.DB.Save(&conf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration updated successfully", "data": conf})
}

func UpdateRateLimit(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		RateLimit  int `json:"rate_limit" binding:"required,min=1"`
		WindowSize int `json:"window_size" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for rate limit update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	appId := c.Param("application_id")

	var conf models.AppConf
	if err := config.DB.Where("application_id = ?", appId).First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App configuration not found"})
		return
	}

	conf.RateLimit = input.RateLimit
	conf.WindowSize = input.WindowSize

	if err := config.DB.Save(&conf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully", "data": conf})
}

func UpdateDetectBot(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		DetectBot bool `json:"detect_bot" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for detect bot update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	appId := c.Param("application_id")

	var conf models.AppConf
	if err := config.DB.Where("application_id = ?", appId).First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App configuration not found"})
		return
	}

	conf.DetectBot = input.DetectBot

	if err := config.DB.Save(&conf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully", "data": conf})
}

func UpdateRemoteLogServer(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		RemoteLogServer string `json:"remote_logServer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging purposes
		log.Printf("Error binding JSON for application update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var conf models.Conf
	if err := config.DB.First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	conf.RemoteLogServer = input.RemoteLogServer

	if err := config.DB.Save(&conf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "configuration updated successfully", "data": conf})
}

func UpdateMaxPosyDataSize(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		MaxPostDataSize float64 `json:"max_post_data_size" `
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for rate limit update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields"})
		return
	}

	appId := c.Param("application_id")

	var conf models.AppConf
	if err := config.DB.Where("application_id = ?", appId).First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App configuration not found"})
		return
	}

	conf.MaxPostDataSize = input.MaxPostDataSize

	if err := config.DB.Save(&conf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully", "data": conf})
}
