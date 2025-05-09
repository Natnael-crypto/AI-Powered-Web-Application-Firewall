package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddApplication(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ApplicationName string `json:"application_name" binding:"required,max=20"`
		Description     string `json:"description" binding:"required,max=200"`
		HostName        string `json:"hostname" binding:"required,max=40"`
		IpAddress       string `json:"ip_address" binding:"required,ip"`
		Port            string `json:"port" binding:"required,max=5"`
		Status          bool   `json:"status"`
		Tls             bool   `json:"tls"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingApp models.Application
	if err := config.DB.Where("application_name = ?", input.ApplicationName).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "application name already exists"})
		return
	}

	if err := config.DB.Where("hostname = ?", input.HostName).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "hostname already exists"})
		return
	}

	application := models.Application{
		ApplicationID:   utils.GenerateUUID(),
		ApplicationName: input.ApplicationName,
		Description:     input.Description,
		HostName:        input.HostName,
		IpAddress:       input.IpAddress,
		Port:            input.Port,
		Status:          input.Status,
		Tls:             input.Tls,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := config.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "failed to create application"})
		return
	}

	newAppConf := models.AppConf{
		ID:              uuid.New().String(),
		ApplicationID:   application.ApplicationID,
		RateLimit:       50,
		WindowSize:      10,
		DetectBot:       false,
		HostName:        application.HostName,
		MaxPostDataSize: 5.0,
		Tls:             true,
	}

	if err := CreateAppConfigLocal(newAppConf); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "failed to create app config"})
	}

	userToApp := models.UserToApplication{
		ID:              utils.GenerateUUID(),
		UserID:          c.GetString("user_id"),
		ApplicationName: application.ApplicationName,
		ApplicationID:   application.ApplicationID,
	}

	if err := config.DB.Create(&userToApp).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "failed to assign user to application"})
		return
	}

	config.Change = true

	c.JSON(http.StatusCreated, gin.H{"message": "application created successfully", "application": application})

}

func GetApplication(c *gin.Context) {
	applicationID := c.Param("application_id")
	appIDs := utils.GetAssignedApplicationIDs(c)

	if !slices.Contains(appIDs, applicationID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	var application models.Application
	if err := config.DB.Where("application_id = ? ", applicationID).Find(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"application": application})

}

func GetAllApplications(c *gin.Context) {
	var applications []models.Application
	if err := config.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}

	var result []map[string]interface{}

	for _, app := range applications {
		appMap := map[string]interface{}{
			"application_id":   app.ApplicationID,
			"application_name": app.ApplicationName,
			"description":      app.Description,
			"hostname":         app.HostName,
			"ip_address":       app.IpAddress,
			"port":             app.Port,
			"status":           app.Status,
			"tls":              app.Tls,
			"created_at":       app.CreatedAt,
			"updated_at":       app.UpdatedAt,
		}

		var appConf models.AppConf
		if err := config.DB.Where("application_id = ?", app.ApplicationID).First(&appConf).Error; err == nil {
			appMap["config"] = appConf
		} else {
			appMap["config"] = gin.H{} // empty config if not found
		}

		result = append(result, appMap)
	}

	c.JSON(http.StatusOK, gin.H{"applications": result})
}

func GetAllApplicationsAdmin(c *gin.Context) {
	appIDs := utils.GetAssignedApplicationIDs(c)

	var applications []models.Application
	if err := config.DB.Where("application_id In ? ", appIDs).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}

	// Merge each application with its config
	var result []map[string]interface{}

	for _, app := range applications {
		appMap := map[string]interface{}{
			"application_id":   app.ApplicationID,
			"application_name": app.ApplicationName,
			"description":      app.Description,
			"hostname":         app.HostName,
			"ip_address":       app.IpAddress,
			"port":             app.Port,
			"status":           app.Status,
			"tls":              app.Tls,
			"created_at":       app.CreatedAt,
			"updated_at":       app.UpdatedAt,
		}

		var appConf models.AppConf
		if err := config.DB.Where("application_id = ?", app.ApplicationID).First(&appConf).Error; err == nil {
			appMap["config"] = appConf
		} else {
			appMap["config"] = gin.H{} // empty config if not found
		}

		result = append(result, appMap)
	}

	c.JSON(http.StatusOK, gin.H{"applications": result})
}

func UpdateApplication(c *gin.Context) {
	applicationID := c.Param("application_id")

	if c.GetString("role") == "super_admin" {
	} else {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var input struct {
		ApplicationName string `json:"application_name" binding:"required,max=20"`
		Description     string `json:"description" binding:"required,max=200"`
		HostName        string `json:"hostname" binding:"required,max=40"`
		IpAddress       string `json:"ip_address" binding:"required,ip"`
		Port            string `json:"port" binding:"required,max=5"`
		Status          bool   `json:"status" binding:"required"`
		Tls             bool   `json:"tls"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for application update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	application.ApplicationName = input.ApplicationName
	application.Description = input.Description
	application.HostName = input.HostName
	application.IpAddress = input.IpAddress
	application.Port = input.Port
	application.Status = input.Status
	application.Tls = input.Tls
	application.UpdatedAt = time.Now()

	if err := config.DB.Model(&application).Where("application_id = ?", applicationID).Updates(map[string]interface{}{
		"application_name": application.ApplicationName,
		"description":      application.Description,
		"hostname":         application.HostName,
		"ip_address":       application.IpAddress,
		"port":             application.Port,
		"status":           application.Status,
		"tls":              application.Tls,
		"updated_at":       application.UpdatedAt,
	}).Error; err != nil {
		log.Printf("Error updating application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update application"})
		return
	}
	config.Change = true
	c.JSON(http.StatusOK, gin.H{"message": "application updated successfully"})
}

func DeleteApplication(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	applicationID := c.Param("application_id")

	if err := config.DB.Where("application_id = ?", applicationID).Delete(&models.Application{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	if err := config.DB.Where("application_id = ?", applicationID).Delete(&models.UserToApplication{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}
