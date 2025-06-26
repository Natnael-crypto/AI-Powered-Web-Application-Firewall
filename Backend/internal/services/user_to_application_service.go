package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUserToApplicationService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		UserID          string `json:"user_id" binding:"required,uuid"`
		ApplicationName string `json:"application_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	_, err := repository.GetUserByID(input.UserID)
	if err != nil {
		return gin.H{"error": "user not found"}, http.StatusNotFound
	}

	application, err := repository.GetApplicationByName(input.ApplicationName)
	if err != nil {
		return gin.H{"error": "application not found"}, http.StatusNotFound
	}

	// Check if the user is already assigned to the application
	if err := repository.CheckUserApplicationAssignment(input.UserID, input.ApplicationName); err == nil {
		return gin.H{"error": "user is already assigned to this application"}, http.StatusConflict
	}

	// Create the assignment
	userToApp := models.UserToApplication{
		ID:              utils.GenerateUUID(),
		UserID:          input.UserID,
		ApplicationName: input.ApplicationName,
		ApplicationID:   application.ApplicationID,
	}

	if err := repository.CreateUserToApplication(userToApp); err != nil {
		return gin.H{"error": "failed to assign user to application"}, http.StatusConflict
	}

	return gin.H{"message": "user assigned to application successfully"}, http.StatusCreated
}

func UpdateUserToApplicationService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		UserID          string `json:"user_id" binding:"required"`
		ApplicationName string `json:"application_name" binding:"required"`
	}

	assignmentID := c.Param("assignment_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	// Fetch the user and application
	_, err := repository.GetUserByID(input.UserID)
	if err != nil {
		return gin.H{"error": "user not found"}, http.StatusNotFound
	}

	_, err = repository.GetApplicationByName(input.ApplicationName)
	if err != nil {
		return gin.H{"error": "application not found"}, http.StatusNotFound
	}

	// Update the assignment
	if err := repository.UpdateUserToApplicationAssignment(assignmentID, input.UserID, input.ApplicationName); err != nil {
		return gin.H{"error": "failed to update user to application assignment"}, http.StatusInternalServerError
	}

	return gin.H{"message": "user to application assignment updated successfully"}, http.StatusOK
}

func GetAllUserToApplicationsService(c *gin.Context) (gin.H, int) {
	assignments, err := repository.GetAllUserToApplications()
	if err != nil {
		return gin.H{"error": "failed to retrieve assignments"}, http.StatusInternalServerError
	}

	return gin.H{"assignments": assignments}, http.StatusOK
}

func DeleteUserToApplicationService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	assignmentID := c.Param("assignment_id")

	// Delete the assignment
	if err := repository.DeleteUserToApplicationAssignment(assignmentID); err != nil {
		return gin.H{"error": "failed to delete assignment"}, http.StatusInternalServerError
	}

	return gin.H{"message": "user to application assignment deleted successfully"}, http.StatusOK
}
