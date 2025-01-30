package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Folder to store certificates
const certStoragePath = "./certs"

// AddCert handles uploading a certificate and key file
func AddCert(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
		return
	}

	applicationID := c.PostForm("application_id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application_id is required"})
		return
	}

	// Check if a certificate already exists for this application
	var existingCert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&existingCert).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Certificate already exists for this application"})
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

	// Define directory path based on application ID
	applicationCertPath := filepath.Join(certStoragePath, applicationID)

	// Ensure storage directory exists
	if err := os.MkdirAll(applicationCertPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application directory"})
		return
	}

	// Define file paths
	certFilePath := filepath.Join(applicationCertPath, certHeader.Filename)
	keyFilePath := filepath.Join(applicationCertPath, keyHeader.Filename)

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
		CertID:        uuid.New().String(),
		CertPath:      certFilePath,
		KeyPath:       keyFilePath,
		ApplicationID: applicationID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := config.DB.Create(&newCert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store certificate in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Certificate uploaded successfully",
		"cert_id":  newCert.CertID,
		"certPath": certFilePath,
		"keyPath":  keyFilePath,
	})
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

// UpdateCert updates the certificate or key file
func UpdateCert(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	applicationID := c.Param("application_id")
	fileType := c.PostForm("type") // Expected values: "cert" or "key"

	if fileType != "cert" && fileType != "key" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type. Must be 'cert' or 'key'"})
		return
	}

	// Fetch the certificate record using ApplicationID
	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Certificate not found for this application"})
		return
	}

	// Get the uploaded file
	file, fileHeader, err := c.Request.FormFile(fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to get %s file", fileType)})
		return
	}
	defer file.Close()

	// Define application-specific folder
	appCertFolder := filepath.Join(certStoragePath, applicationID)
	if err := os.MkdirAll(appCertFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application cert folder"})
		return
	}

	// Define new file path
	newFilePath := filepath.Join(appCertFolder, fileHeader.Filename)

	// Save the new file
	outFile, err := os.Create(newFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save %s file", fileType)})
		return
	}
	defer outFile.Close()
	io.Copy(outFile, file)

	// Remove the old file
	var oldFilePath string
	if fileType == "cert" {
		oldFilePath = cert.CertPath
		cert.CertPath = newFilePath
	} else {
		oldFilePath = cert.KeyPath
		cert.KeyPath = newFilePath
	}

	// Delete old file if it exists
	if oldFilePath != "" {
		if err := os.Remove(oldFilePath); err != nil {
			fmt.Printf("Warning: Could not delete old file %s\n", oldFilePath)
		}
	}

	// Update the database record
	if err := config.DB.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update certificate path in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("%s file updated successfully", fileType),
		"new_path": newFilePath,
		"cert_id":  cert.CertID,
		"app_id":   cert.ApplicationID,
	})
}

// DeleteCert deletes a certificate
func DeleteCert(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

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
