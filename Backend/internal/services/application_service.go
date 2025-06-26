package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"slices"
	"time"
)

func AddApplicationService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input models.ApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	if repository.ApplicationExistsByName(input.ApplicationName) {
		return gin.H{"error": "application name already exists"}, http.StatusConflict
	}
	if repository.ApplicationExistsByHost(input.HostName) {
		return gin.H{"error": "hostname already exists"}, http.StatusConflict
	}

	app := models.Application{
		ApplicationID:   utils.GenerateUUID(),
		ApplicationName: input.ApplicationName,
		Description:     input.Description,
		HostName:        input.HostName,
		IpAddress:       input.IpAddress,
		Port:            input.Port,
		Status:          *input.Status,
		Tls:             *input.Tls,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := repository.CreateApplication(app); err != nil {
		return gin.H{"error": "failed to create application"}, http.StatusConflict
	}

	appConf := models.AppConf{
		ID:              uuid.New().String(),
		ApplicationID:   app.ApplicationID,
		RateLimit:       50,
		WindowSize:      10,
		DetectBot:       false,
		HostName:        app.HostName,
		MaxPostDataSize: 5.0,
		Tls:             false,
	}
	if err := repository.CreateAppConfig(appConf); err != nil {
		return gin.H{"error": "failed to create app config"}, http.StatusConflict
	}

	userToApp := models.UserToApplication{
		ID:              utils.GenerateUUID(),
		UserID:          c.GetString("user_id"),
		ApplicationID:   app.ApplicationID,
		ApplicationName: app.ApplicationName,
	}
	if err := repository.CreateUserToApp(userToApp); err != nil {
		return gin.H{"error": "failed to assign user to application"}, http.StatusConflict
	}

	if !app.Tls {
		repository.CreateEmptyCertificate(app.ApplicationID)
	}

	config.Change = true
	return gin.H{"message": "application created successfully", "application": app}, http.StatusCreated
}

func GetApplicationService(c *gin.Context) (gin.H, int) {
	appID := c.Param("application_id")
	if !slices.Contains(utils.GetAssignedApplicationIDs(c), appID) {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}
	application, err := repository.GetApplicationByID(appID)
	if err != nil {
		return gin.H{"error": "failed to fetch application"}, http.StatusInternalServerError
	}
	return gin.H{"application": application}, http.StatusOK
}

func GetAllApplicationsService(c *gin.Context) (gin.H, int) {
	data, err := repository.FetchAllApplicationsWithConfig()
	if err != nil {
		return gin.H{"error": "failed to fetch applications"}, http.StatusInternalServerError
	}
	return gin.H{"applications": data}, http.StatusOK
}

func GetAdminApplicationsService(c *gin.Context) (gin.H, int) {
	appIDs := utils.GetAssignedApplicationIDs(c)
	data, err := repository.FetchApplicationsByIDsWithConfig(appIDs)
	if err != nil {
		return gin.H{"error": "failed to fetch applications"}, http.StatusInternalServerError
	}
	return gin.H{"applications": data}, http.StatusOK
}

func UpdateApplicationService(c *gin.Context) (gin.H, int) {
	appID := c.Param("application_id")
	if c.GetString("role") != "super_admin" && !utils.HasAccessToApplication(utils.GetAssignedApplicationIDs(c), appID) {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input models.ApplicationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid input: %v", err)
		return gin.H{"error": "invalid input format"}, http.StatusBadRequest
	}

	err := repository.UpdateApplication(appID, input)
	if err != nil {
		log.Printf("Update error: %v", err)
		return gin.H{"error": "failed to update application"}, http.StatusInternalServerError
	}

	config.Change = true
	return gin.H{"message": "application updated successfully"}, http.StatusOK
}

func DeleteApplicationService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}
	appID := c.Param("application_id")
	if err := repository.DeleteApplicationAndLinks(appID); err != nil {
		return gin.H{"error": "application not found"}, http.StatusNotFound
	}
	config.Change = true
	return gin.H{"message": "application deleted successfully"}, http.StatusOK
}
