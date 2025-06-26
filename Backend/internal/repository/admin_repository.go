package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func AdminExists(username string) bool {
	var existing models.User
	err := config.DB.Where("username = ?", username).First(&existing).Error
	return err == nil
}

func CreateAdmin(admin models.User) error {
	return config.DB.Create(&admin).Error
}

func GetAdminByUsername(username string) (models.User, error) {
	var admin models.User
	err := config.DB.Where("username = ?", username).First(&admin).Error
	return admin, err
}

func GetAdminByID(id string) (models.User, error) {
	var admin models.User
	err := config.DB.Where("user_id = ?", id).First(&admin).Error
	return admin, err
}

func GetAllAdmins() ([]models.User, error) {
	var admins []models.User
	err := config.DB.Find(&admins).Error
	return admins, err
}

func DeleteAdminByUsername(username string) error {
	return config.DB.Where("username = ?", username).Delete(&models.User{}).Error
}

func UpdateAdminStatus(username string, status string) error {
	return config.DB.Model(&models.User{}).Where("username = ?", username).Update("status", status).Error
}
