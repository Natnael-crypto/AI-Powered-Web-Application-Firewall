package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateConfig handles the creation of a new config entry
func CreateConfig(c *gin.Context) {

	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ListeningPort   string `json:"listening_port" binding:"required"`
		RemoteLogServer string `json:"remote_logServer" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a config entry already exists
	var conf models.Conf
	if err := config.DB.First(&conf).Error; err == nil {
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
func UpdateConfig(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ListeningPort   string `json:"listening_port" binding:"required"`
		RemoteLogServer string `json:remote_logServer binding:"required"`
	}

	configID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var conf models.Conf
	if err := config.DB.Where("id = ?", configID).First(&conf).Error; err != nil {
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

// GetConfig retrieves a configuration by ID
func GetConfig(c *gin.Context) {
	var conf models.Conf
	if err := config.DB.First(&conf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conf})
}
