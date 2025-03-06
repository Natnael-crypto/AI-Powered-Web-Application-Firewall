package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ApplicationName string `json:"application_name" binding:"required,max=20"`
		Description     string `json:"description" binding:"required,max=200"`
		HostName        string `json:"hostname" binding:"required,max=40"`
		IpAddress       string `json:"ip_address" binding:"required"`
		Port            string `json:"port" binding:"required,max=5"`
		Status          bool   `json:"status"`
		Tls             bool   `json:"tls"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the application name or hostname already exists
	var existingApp models.Application
	if err := config.DB.Where("application_name = ?", input.ApplicationName).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "application name already exists"})
		return
	}

	if err := config.DB.Where("hostname = ?", input.HostName).First(&existingApp).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "hostname already exists"})
		return
	}

	// Create the new application
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

	// Save the application to the PostgreSQL database
	if err := config.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "application created successfully"})
}

// GetApplication retrieves a specific application by ID
func GetApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	applicationID := c.Param("application_id")
	var application models.Application

	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": application})
}

// GetAllApplications retrieves all applications
func GetAllApplications(c *gin.Context) {

	var applications []models.Application
	if err := config.DB.Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// UpdateApplication updates an existing application
func UpdateApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		ApplicationName string `json:"application_name" binding:"required,max=20"`
		Description     string `json:"description" binding:"required,max=200"`
		HostName        string `json:"hostname" binding:"required,max=40"`
		IpAddress       string `json:"ip_address" binding:"required,max=15"`
		Port            string `json:"port" binding:"required,max=5"`
		Status          bool   `json:"status" binding:"required"`
		Tls             bool   `json:"tls"`
	}

	applicationID := c.Param("application_id")

	// Check if the request body is empty
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging purposes
		log.Printf("Error binding JSON for application update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Update fields in the application struct
	application.ApplicationName = input.ApplicationName
	application.Description = input.Description
	application.HostName = input.HostName
	application.IpAddress = input.IpAddress
	application.Port = input.Port
	application.Status = input.Status
	application.Tls = input.Tls
	application.UpdatedAt = time.Now()

	// Explicitly update the application using WHERE condition
	if err := config.DB.Model(&application).Where("application_id = ?", applicationID).Updates(map[string]interface{}{
		"application_name": application.ApplicationName,
		"description":      application.Description,
		"host_name":        application.HostName,
		"ip_address":       application.IpAddress,
		"port":             application.Port,
		"status":           application.Status,
		"tls":              application.Tls,
		"updated_at":       application.UpdatedAt,
	}).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error updating application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application updated successfully"})
}

// DeleteApplication deletes an application by ID
func DeleteApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	applicationID := c.Param("application_id")

	if err := config.DB.Where("application_id = ?", applicationID).Delete(&models.Application{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}
