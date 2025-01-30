package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Folder to store certificates
const certStoragePath = "./certs"

// AddCert handles uploading a certificate and key file
func AddCert(c *gin.Context) {

	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	applicationID := c.PostForm("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application_id is required"})
		return
	}

	// Parse files
	certFile, certHeader, err := c.Request.FormFile("cert")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get cert file"})
		return
	}
	defer certFile.Close()

	keyFile, keyHeader, err := c.Request.FormFile("key")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get key file"})
		return
	}
	defer keyFile.Close()

	// Ensure storage directory exists
	if err := os.MkdirAll(certStoragePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
		return
	}

	// Define file paths
	certFilePath := filepath.Join(certStoragePath, certHeader.Filename)
	keyFilePath := filepath.Join(certStoragePath, keyHeader.Filename)

	// Save cert file
	outCert, err := os.Create(certFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cert file"})
		return
	}
	defer outCert.Close()
	io.Copy(outCert, certFile)

	// Save key file
	outKey, err := os.Create(keyFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save key file"})
		return
	}
	defer outKey.Close()
	io.Copy(outKey, keyFile)

	// Save to DB
	newCert := models.Cert{
		CertPath:      certFilePath,
		KeyPath:       keyFilePath,
		ApplicationID: applicationID,
	}
	if err := config.DB.Create(&newCert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store certificate in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Certificate uploaded successfully", "cert_id": newCert.CertID})
}

// GetCert retrieves a certificate or key file
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

	filePath := cert.CertPath
	if fileType == "key" {
		filePath = cert.KeyPath
	}

	c.File(filePath)
}

// UpdateCert updates certificate paths
func UpdateCert(c *gin.Context) {
	certID := c.Param("cert_id")
	var cert models.Cert

	if err := config.DB.First(&cert, certID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := config.DB.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate updated successfully"})
}

// DeleteCert deletes a certificate
func DeleteCert(c *gin.Context) {
	certID := c.Param("cert_id")
	var cert models.Cert

	if err := config.DB.First(&cert, certID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found"})
		return
	}

	// Remove files from storage
	if err := os.Remove(cert.CertPath); err != nil {
		fmt.Println("Warning: Could not delete cert file")
	}
	if err := os.Remove(cert.KeyPath); err != nil {
		fmt.Println("Warning: Could not delete key file")
	}

	// Delete from DB
	if err := config.DB.Delete(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete certificate from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Certificate deleted successfully"})
}
