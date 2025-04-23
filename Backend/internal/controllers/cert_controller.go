package controllers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCert(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
		return
	}

	applicationID := c.Param("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application_id is required"})
		return
	}

	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var existingCert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&existingCert).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Certificate already exists for this application"})
		return
	}

	certFile, _, err := c.Request.FormFile("cert")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get cert file"})
		return
	}
	defer certFile.Close()

	keyFile, _, err := c.Request.FormFile("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get key file"})
		return
	}
	defer keyFile.Close()

	certContent, err := io.ReadAll(certFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read cert file"})
		return
	}

	keyContent, err := io.ReadAll(keyFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read key file"})
		return
	}

	newCert := models.Cert{
		CertID:        uuid.New().String(),
		Cert:          certContent,
		Key:           keyContent,
		ApplicationID: applicationID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&newCert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store certificate in database"})
		return
	}

	config.Change = true

	c.JSON(http.StatusCreated, gin.H{
		"message": "Certificate uploaded successfully",
		"cert_id": newCert.CertID,
	})
}

func GetCert(c *gin.Context) {
	applicationID := c.Query("application_id")
	fileType := c.Query("type")

	if applicationID == "" || (fileType != "cert" && fileType != "key") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	var file []byte
	if fileType == "key" {
		file = cert.Key
	} else {
		file = cert.Cert
	}

	c.Data(http.StatusOK, "application/octet-stream", file)
}

func UpdateCert(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}
	applicationID := c.Param("application_id")
	fileType := c.PostForm("type")

	if fileType != "cert" && fileType != "key" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type. Must be 'cert' or 'key'"})
		return
	}
	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found for this application"})
		return
	}

	file, _, err := c.Request.FormFile(fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to get %s file", fileType)})
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to read %s file", fileType)})
		return
	}

	if fileType == "cert" {
		cert.Cert = fileContent
	} else {
		cert.Key = fileContent
	}

	if err := config.DB.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate in database"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s file updated successfully", fileType),
		"cert_id": cert.CertID,
		"app_id":  cert.ApplicationID,
	})
}

func DeleteCert(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	certID := c.Param("application_id")
	var cert models.Cert

	if err := config.DB.Where("application_id = ?", certID).First(&cert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	if err := config.DB.Delete(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete certificate from database"})
		return
	}

	config.Change = true

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
