package controllers

import (
	"net/http"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AddUserToApplication assigns a user to an application
func AddUserToApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		UserID        string `json:"user_id" binding:"required"`
		ApplicationName string `json:"application_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if application exists
	var application models.Application
	if err := config.DB.Where("application_Name = ?", input.ApplicationName).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Check if the user is already assigned to the application
	var existingAssignment models.UserToApplication
	if err := config.DB.Where("user_id = ? AND application_name = ?", input.UserID, input.ApplicationName).First(&existingAssignment).Error; err == nil {
		// If assignment exists, return an error
		c.JSON(http.StatusConflict, gin.H{"error": "user is already assigned to this application"})
		return
	}

	// Create the assignment
	userToApp := models.UserToApplication{
		ID:            utils.GenerateUUID(),
		UserID:        input.UserID,
		ApplicationName: input.ApplicationName,
	}

	if err := config.DB.Create(&userToApp).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "failed to assign user to application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user assigned to application successfully"})
}

// UpdateUserToApplication updates the assignment of a user to an application
func UpdateUserToApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		UserID        string `json:"user_id" binding:"required"`
		ApplicationName string `json:"application_name" binding:"required"`
	}

	assignmentID := c.Param("assignment_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.Where("user_id = ?", input.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if application exists
	var application models.Application
	if err := config.DB.Where("application_name = ?", input.ApplicationName).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	// Update the assignment
	if err := config.DB.Model(&models.UserToApplication{}).Where("id = ?", assignmentID).Updates(map[string]interface{}{
		"user_id":        input.UserID,
		"application_name": input.ApplicationName,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user to application assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user to application assignment updated successfully"})
}

// GetAllUserToApplications retrieves all user to application assignments
func GetAllUserToApplications(c *gin.Context) {
	var assignments []models.UserToApplication
	if err := config.DB.Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve assignments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"assignments": assignments})
}

// DeleteUserToApplication deletes the assignment of a user to an application
func DeleteUserToApplication(c *gin.Context) {
	// Check if the user is a super admin
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	assignmentID := c.Param("assignment_id")

	// Check if the assignment exists
	var assignment models.UserToApplication
	if err := config.DB.Where("id = ?", assignmentID).First(&assignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assignment not found"})
		return
	}

	// Delete the assignment
	if err := config.DB.Delete(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user to application assignment deleted successfully"})
}
