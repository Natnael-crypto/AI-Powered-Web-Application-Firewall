package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

// GetSystemEmail retrieves the system email configuration
func GetSystemEmail() (models.SystemEmail, error) {
	var sysEmail models.SystemEmail
	if err := config.DB.First(&sysEmail).Error; err != nil {
		return sysEmail, err
	}
	return sysEmail, nil
}

// CreateSystemEmail inserts a new system email configuration into the database
func CreateSystemEmail(sysEmail models.SystemEmail) error {
	return config.DB.Create(&sysEmail).Error
}

// UpdateSystemEmail updates the system email configuration in the database
func UpdateSystemEmail(sysEmail models.SystemEmail) error {
	return config.DB.Save(&sysEmail).Error
}
