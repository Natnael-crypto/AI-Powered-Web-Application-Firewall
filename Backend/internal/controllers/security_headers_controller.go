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

func AddSecurityHeader(c *gin.Context) {
	userID := c.GetString("user_id")

	var input struct {
		HeaderName     string   `json:"header_name" binding:"required,max=50"`
		HeaderValue    string   `json:"header_value" binding:"required,max=500"`
		ApplicationIDs []string `json:"application_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appIds := utils.GetAssignedApplicationIDs(c)

	for _, id := range input.ApplicationIDs {
		if !slices.Contains(appIds, id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	securityHeader := models.SecurityHeader{
		ID:          uuid.New().String(),
		HeaderName:  input.HeaderName,
		HeaderValue: input.HeaderValue,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, id := range input.ApplicationIDs {
		applicationToSecurity := models.ApplicationSecurityHeader{
			ID:               uuid.New().String(),
			ApplicationID:    id,
			SecurityHeaderID: securityHeader.ID,
		}

		if err := config.DB.Create(&applicationToSecurity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create security header to application relationship"})
			return
		}

	}

	if err := config.DB.Create(&securityHeader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusCreated, gin.H{"message": "security header added successfully"})
}

func GetSecurityHeaders(c *gin.Context) {
	applicationID := c.Param("application_id")

	var applicationToSecurityHeaders []models.ApplicationSecurityHeader

	if err := config.DB.Where("application_id = ?", applicationID).Find(&applicationToSecurityHeaders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch application to security headers relation"})
		return
	}

	appSec := make([]string, 0, len(applicationToSecurityHeaders))
	for _, mapping := range applicationToSecurityHeaders {
		appSec = append(appSec, mapping.SecurityHeaderID)
	}

	var securityHeaders []models.SecurityHeader
	query := config.DB.Model(&models.SecurityHeader{})

	if appSec == nil {
		c.JSON(http.StatusOK, gin.H{"security_headers": securityHeaders})
		return
	}

	if err := query.Where("id In ?", appSec).Find(&securityHeaders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch security headers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"security_headers": securityHeaders})
}

func GetSecurityHeadersAdmin(c *gin.Context) {

	appIDs := utils.GetAssignedApplicationIDs(c)

	var applicationToSecurityHeaders []models.ApplicationSecurityHeader

	if err := config.DB.Where("application_id In ?", appIDs).Find(&applicationToSecurityHeaders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch application to security headers relation"})
		return
	}

	appSec := make([]string, 0, len(applicationToSecurityHeaders))
	for _, mapping := range applicationToSecurityHeaders {
		appSec = append(appSec, mapping.SecurityHeaderID)
	}

	var securityHeaders []models.SecurityHeader
	query := config.DB.Model(&models.SecurityHeader{})

	if err := query.Where("id In ?", appSec).Find(&securityHeaders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch security headers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"security_headers": securityHeaders})
}

func UpdateSecurityHeader(c *gin.Context) {
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	var input struct {
		HeaderName  string `json:"header_name" binding:"required,max=50"`
		HeaderValue string `json:"header_value" binding:"required,max=500"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for security header update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	var securityHeader models.SecurityHeader
	if err := config.DB.Where("id = ? And created_by =?", headerID, userID).First(&securityHeader).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security header not found"})
		return
	}

	if input.HeaderName != "" {
		securityHeader.HeaderName = input.HeaderName
	}
	if input.HeaderValue != "" {
		securityHeader.HeaderValue = input.HeaderValue
	}

	securityHeader.UpdatedAt = time.Now()

	if err := config.DB.Save(&securityHeader).Where("id = ?", headerID).Error; err != nil {
		log.Printf("Error updating security header: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "security header updated successfully"})
}

func DeleteSecurityHeader(c *gin.Context) {
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	var securityHeader models.SecurityHeader
	if err := config.DB.Where("id = ? And created_by =?", headerID, userID).First(&securityHeader).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "security header not found"})
		return
	}

	if err := config.DB.Delete(&securityHeader).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete security header"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "security header deleted successfully"})
}
