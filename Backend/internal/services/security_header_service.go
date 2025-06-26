package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddSecurityHeaderService(c *gin.Context) (gin.H, int) {
	userID := c.GetString("user_id")
	var input struct {
		HeaderName     string   `json:"header_name" binding:"required,max=50"`
		HeaderValue    string   `json:"header_value" binding:"required,max=500"`
		ApplicationIDs []string `json:"application_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	appIDs := utils.GetAssignedApplicationIDs(c)
	for _, id := range input.ApplicationIDs {
		if !slices.Contains(appIDs, id) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	header := models.SecurityHeader{
		ID:            uuid.New().String(),
		HeaderName:    input.HeaderName,
		HeaderValue:   input.HeaderValue,
		ApplicationID: input.ApplicationIDs[0],
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := repository.CreateSecurityHeader(header); err != nil {
		return gin.H{"error": "failed to create security header"}, http.StatusInternalServerError
	}

	config.Change = true
	return gin.H{"message": "security header added successfully"}, http.StatusCreated
}

func GetSecurityHeadersService(c *gin.Context) (gin.H, int) {
	appID := c.Param("application_id")
	headers, err := repository.GetSecurityHeadersByAppID(appID)
	if err != nil {
		return gin.H{"error": "failed to fetch security headers"}, http.StatusInternalServerError
	}
	return gin.H{"security_headers": headers}, http.StatusOK
}

func GetSecurityHeadersAdminService(c *gin.Context) (gin.H, int) {
	appIDs := utils.GetAssignedApplicationIDs(c)
	headers, err := repository.GetSecurityHeadersByAppIDs(appIDs)
	if err != nil {
		return gin.H{"error": "failed to fetch security headers"}, http.StatusInternalServerError
	}
	return gin.H{"security_headers": headers}, http.StatusOK
}

func UpdateSecurityHeaderService(c *gin.Context) (gin.H, int) {
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	var input struct {
		HeaderName  string `json:"header_name" binding:"required,max=50"`
		HeaderValue string `json:"header_value" binding:"required,max=500"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON for update: %v", err)
		return gin.H{"error": "Invalid input format"}, http.StatusBadRequest
	}

	header, err := repository.GetSecurityHeaderByIDAndUser(headerID, userID)
	if err != nil {
		return gin.H{"error": "security header not found"}, http.StatusNotFound
	}

	header.HeaderName = input.HeaderName
	header.HeaderValue = input.HeaderValue
	header.UpdatedAt = time.Now()

	if err := repository.UpdateSecurityHeader(header); err != nil {
		log.Printf("Update error: %v", err)
		return gin.H{"error": "failed to update security header"}, http.StatusInternalServerError
	}

	config.Change = true
	return gin.H{"message": "security header updated successfully"}, http.StatusOK
}

func DeleteSecurityHeaderService(c *gin.Context) (gin.H, int) {
	userID := c.GetString("user_id")
	headerID := c.Param("header_id")

	header, err := repository.GetSecurityHeaderByIDAndUser(headerID, userID)
	if err != nil {
		return gin.H{"error": "insufficient privileges"}, http.StatusNotFound
	}

	if err := repository.DeleteSecurityHeader(header); err != nil {
		return gin.H{"error": "failed to delete security header"}, http.StatusInternalServerError
	}

	config.Change = true
	return gin.H{"message": "security header deleted successfully"}, http.StatusOK
}
