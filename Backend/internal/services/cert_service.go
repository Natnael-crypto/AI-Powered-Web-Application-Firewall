package services

import (
	"backend/internal/repository"
	"backend/internal/utils"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func validateCertificate(certPEM []byte) error {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return errors.New("failed to parse PEM block from certificate")
	}
	_, err := x509.ParseCertificate(block.Bytes)
	return err
}

func validatePrivateKey(keyPEM []byte) error {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return errors.New("failed to parse PEM block from key")
	}

	if _, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return nil
	}
	if _, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return nil
	}
	if _, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return nil
	}

	return errors.New("unsupported or invalid private key format")
}

func AddCertService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	if applicationID == "" {
		return gin.H{"error": "application_id is required"}, http.StatusBadRequest
	}

	// Permission check
	if c.GetString("role") != "super_admin" {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	// Validate and get certificate & key
	certFile, keyFile, err := getFiles(c)
	if err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	// Validate the certificate and private key
	if err := validateCertificate(certFile); err != nil {
		return gin.H{"error": "Invalid certificate file: " + err.Error()}, http.StatusBadRequest
	}

	if err := validatePrivateKey(keyFile); err != nil {
		return gin.H{"error": "Invalid key file: " + err.Error()}, http.StatusBadRequest
	}

	// Update database
	return repository.AddCertToDatabase(applicationID, certFile, keyFile)
}

func GetCertService(c *gin.Context) (gin.H, int) {
	applicationID := c.Query("application_id")
	fileType := c.Query("type")
	if applicationID == "" || (fileType != "cert" && fileType != "key") {
		return gin.H{"error": "Invalid request parameters"}, http.StatusBadRequest
	}
	return repository.GetCertFromDatabase(applicationID, fileType)
}

func GetCertAdminService(c *gin.Context) (gin.H, int) {
	applicationID := c.Query("application_id")
	fileType := c.Query("type")
	if applicationID == "" || (fileType != "cert" && fileType != "key") {
		return gin.H{"error": "Invalid request parameters"}, http.StatusBadRequest
	}

	// Permission check
	if c.GetString("role") != "super_admin" {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	return repository.GetCertFromDatabase(applicationID, fileType)
}

func UpdateCertService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")
	fileType := c.PostForm("type")

	// Validate file type
	if fileType != "cert" && fileType != "key" {
		return gin.H{"error": "Invalid type. Must be 'cert' or 'key'"}, http.StatusBadRequest
	}

	// Permission check
	if c.GetString("role") != "super_admin" {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	// Get file and validate
	file, err := getFile(c, fileType)
	if err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	// Update the certificate in DB
	return repository.UpdateCertInDatabase(applicationID, file, fileType)
}

func DeleteCertService(c *gin.Context) (gin.H, int) {
	applicationID := c.Param("application_id")

	// Permission check
	if c.GetString("role") != "super_admin" {
		appIds := utils.GetAssignedApplicationIDs(c)
		if !utils.HasAccessToApplication(appIds, applicationID) {
			return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
		}
	}

	return repository.DeleteCertFromDatabase(applicationID)
}

// Helper functions
func getFiles(c *gin.Context) ([]byte, []byte, error) {
	certFile, _, err := c.Request.FormFile("cert")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cert file")
	}
	defer certFile.Close()

	keyFile, _, err := c.Request.FormFile("key")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get key file")
	}
	defer keyFile.Close()

	certContent, err := io.ReadAll(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cert file")
	}

	keyContent, err := io.ReadAll(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read key file")
	}

	return certContent, keyContent, nil
}

func getFile(c *gin.Context, fileType string) ([]byte, error) {
	file, _, err := c.Request.FormFile(fileType)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s file", fileType)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file", fileType)
	}

	return fileContent, nil
}
