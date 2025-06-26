package repository

import (
	"backend/internal/config"
	"backend/internal/models"
	"errors"
)

// GetUserByID fetches a user by their user ID
// func GetUserByID(userID string) (models.User, error) {
// 	var user models.User
// 	if err := config.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
// 		return user, err
// 	}
// 	return user, nil
// }

// GetApplicationByName fetches an application by its name
func GetApplicationByName(applicationName string) (models.Application, error) {
	var application models.Application
	if err := config.DB.Where("application_name = ?", applicationName).First(&application).Error; err != nil {
		return application, err
	}
	return application, nil
}

// CheckUserApplicationAssignment checks if a user is already assigned to an application
func CheckUserApplicationAssignment(userID, applicationName string) error {
	var assignment models.UserToApplication
	if err := config.DB.Where("user_id = ? AND application_name = ?", userID, applicationName).First(&assignment).Error; err == nil {
		return nil
	}
	return errors.New("user is not assigned to the application")
}

// CreateUserToApplication creates a user-to-application assignment
func CreateUserToApplication(userToApp models.UserToApplication) error {
	return config.DB.Create(&userToApp).Error
}

// UpdateUserToApplicationAssignment updates the user-to-application assignment
func UpdateUserToApplicationAssignment(assignmentID, userID, applicationName string) error {
	return config.DB.Model(&models.UserToApplication{}).Where("id = ?", assignmentID).Updates(map[string]interface{}{
		"user_id":          userID,
		"application_name": applicationName,
	}).Error
}

// GetAllUserToApplications retrieves all user-to-application assignments
func GetAllUserToApplications() ([]models.UserToApplication, error) {
	var assignments []models.UserToApplication
	if err := config.DB.Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

// DeleteUserToApplicationAssignment deletes a user-to-application assignment by ID
func DeleteUserToApplicationAssignment(assignmentID string) error {
	return config.DB.Delete(&models.UserToApplication{}, assignmentID).Error
}
