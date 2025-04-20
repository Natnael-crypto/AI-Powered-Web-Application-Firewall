package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddSecurityHeader adds a new security header
func AddSecurityHeader(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")

	var input struct {
		HeaderName    string `json:"header_name" binding:"required,max=50"`
		HeaderValue   string `json:"header_value" binding:"required,max=500"`
		ApplicationID string `json:"application_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If the user is not a super admin and an application ID is provided, verify assignment
	if userRole != "super_admin" {
		var userToApp models.UserToApplication
		if err := config.DB.Where("user_id = ? AND application_id = ?", userID, input.ApplicationID).First(&userToApp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	// Check if the header name already exists for the given application ID
	var existingHeader models.SecurityHeader
	if err := config.DB.
		Where("header_name = ? AND application_id = ?", input.HeaderName, input.ApplicationID).
		First(&existingHeader).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "security header already exists for this application"})
		return
	}

	securityHeader := models.SecurityHeader{
		ID:            uuid.New().String(),
		ApplicationID: input.ApplicationID,
		HeaderName:    input.HeaderName,
		HeaderValue:   input.HeaderValue,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&securityHeader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusCreated, gin.H{"message": "security header added successfully"})
}

// GetSecurityHeaders retrieves security headers based on user role and application
func GetSecurityHeaders(c *gin.Context) {
	applicationID := c.Param("application_id")

	var securityHeaders []models.SecurityHeader
	query := config.DB.Model(&models.SecurityHeader{})

	// Filter by application ID if provided
	if applicationID != "" {
		query = query.Where("application_id = ?", applicationID)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "parameter application_id is empty"})
		return
	}

	if err := query.Find(&securityHeaders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch security headers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"security_headers": securityHeaders})
}

// UpdateSecurityHeader updates an existing security header
func UpdateSecurityHeader(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	var input struct {
		HeaderName  string `json:"header_name" binding:"required ,max=50"`
		HeaderValue string `json:"header_value" binding:"required ,max=500"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for security header update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var securityHeader models.SecurityHeader
	if err := config.DB.Where("id = ?", headerID).First(&securityHeader).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security header not found"})
		return
	}

	if userRole != "super_admin" {
		var userToApp models.UserToApplication
		if err := config.DB.Where("user_id = ? AND application_id = ?", userID, securityHeader.ApplicationID).First(&userToApp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	// Update security header fields
	if input.HeaderName != "" {
		securityHeader.HeaderName = input.HeaderName
	}
	if input.HeaderValue != "" {
		securityHeader.HeaderValue = input.HeaderValue
	}

	securityHeader.UpdatedAt = time.Now()

	// Explicitly update only modified fields
	if err := config.DB.Model(&securityHeader).Where("id = ?", headerID).Updates(map[string]interface{}{
		"header_name":  securityHeader.HeaderName,
		"header_value": securityHeader.HeaderValue,
		"updated_at":   securityHeader.UpdatedAt,
	}).Error; err != nil {
		log.Printf("Error updating security header: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "security header updated successfully"})
}

// DeleteSecurityHeader removes a security header
func DeleteSecurityHeader(c *gin.Context) {
	userRole := c.GetString("role")
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	var securityHeader models.SecurityHeader
	if err := config.DB.Where("id = ?", headerID).First(&securityHeader).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security header not found"})
		return
	}

	if userRole != "super_admin" {
		var userToApp models.UserToApplication
		if err := config.DB.Where("user_id = ? AND application_id = ?", userID, securityHeader.ApplicationID).First(&userToApp).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	if err := config.DB.Delete(&securityHeader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "security header deleted successfully"})
}
