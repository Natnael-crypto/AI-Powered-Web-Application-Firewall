package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func CreateSecurityHeader(header models.SecurityHeader) error {
	return config.DB.Create(&header).Error
}

func GetSecurityHeadersByAppID(appID string) ([]models.SecurityHeader, error) {
	var headers []models.SecurityHeader
	err := config.DB.Where("application_id = ?", appID).Find(&headers).Error
	return headers, err
}

func GetSecurityHeadersByAppIDs(appIDs []string) ([]models.SecurityHeader, error) {
	var headers []models.SecurityHeader
	err := config.DB.Where("application_id IN ?", appIDs).Find(&headers).Error
	return headers, err
}

func GetSecurityHeaderByIDAndUser(id, userID string) (models.SecurityHeader, error) {
	var header models.SecurityHeader
	err := config.DB.Where("id = ? AND created_by = ?", id, userID).First(&header).Error
	return header, err
}

func UpdateSecurityHeader(header models.SecurityHeader) error {
	return config.DB.Save(&header).Error
}

func DeleteSecurityHeader(header models.SecurityHeader) error {
	return config.DB.Delete(&header).Error
}
