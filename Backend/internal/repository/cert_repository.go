package repository

import (
	"backend/internal/config"
	"backend/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddCertToDatabase(applicationID string, certContent, keyContent []byte) (gin.H, int) {
	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		return gin.H{"error": "Application not found"}, http.StatusNotFound
	}

	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		return gin.H{"error": "Certificate doesn't exist for this application"}, http.StatusConflict
	}

	cert.Cert = certContent
	cert.Key = keyContent
	cert.UpdatedAt = time.Now()

	if err := config.DB.Save(&cert).Error; err != nil {
		return gin.H{"error": "Failed to store certificate in database"}, http.StatusInternalServerError
	}

	application.Tls = true
	if err := config.DB.Save(&application).Error; err != nil {
		return gin.H{"error": "Failed to update application TLS status"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Certificate uploaded successfully", "cert_id": cert.CertID}, http.StatusCreated
}

func GetCertFromDatabase(applicationID, fileType string) (gin.H, int) {
	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		return gin.H{"error": "Certificate not found"}, http.StatusNotFound
	}

	var file []byte
	if fileType == "key" {
		file = cert.Key
	} else {
		file = cert.Cert
	}

	return gin.H{"file": file}, http.StatusOK
}

func UpdateCertInDatabase(applicationID string, fileContent []byte, fileType string) (gin.H, int) {
	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		return gin.H{"error": "Certificate not found for this application"}, http.StatusNotFound
	}

	if fileType == "cert" {
		cert.Cert = fileContent
	} else {
		cert.Key = fileContent
	}

	if err := config.DB.Save(&cert).Error; err != nil {
		return gin.H{"error": "Failed to update certificate in database"}, http.StatusInternalServerError
	}

	return gin.H{"message": fmt.Sprintf("%s file updated successfully", fileType)}, http.StatusOK
}

func DeleteCertFromDatabase(applicationID string) (gin.H, int) {
	var cert models.Cert
	if err := config.DB.Where("application_id = ?", applicationID).First(&cert).Error; err != nil {
		return gin.H{"error": "Certificate not found"}, http.StatusNotFound
	}

	if err := config.DB.Delete(&cert).Error; err != nil {
		return gin.H{"error": "Failed to delete certificate from database"}, http.StatusInternalServerError
	}

	var application models.Application
	if err := config.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		return gin.H{"error": "Application not found"}, http.StatusNotFound
	}

	application.Tls = false
	if err := config.DB.Save(&application).Error; err != nil {
		return gin.H{"error": "Failed to update application TLS status"}, http.StatusInternalServerError
	}

	return gin.H{"message": "Certificate deleted successfully"}, http.StatusOK
}
