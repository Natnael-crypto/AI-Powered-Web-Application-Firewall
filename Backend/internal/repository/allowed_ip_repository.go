package repository

import (
	"backend/internal/config"
	"backend/internal/models"
	"errors"
)

// CreateAllowedIp inserts a new AllowedIp into the database
func CreateAllowedIp(allowedIp models.AllowedIp) error {
	return config.DB.Create(&allowedIp).Error
}

// GetAllowedIpByID fetches an AllowedIp by its ID
func GetAllowedIpByID(id string) (models.AllowedIp, error) {
	var ip models.AllowedIp
	if err := config.DB.Where("id = ?", id).First(&ip).Error; err != nil {
		return models.AllowedIp{}, err
	}
	return ip, nil
}

// UpdateAllowedIp updates an existing AllowedIp
func UpdateAllowedIp(ip models.AllowedIp) error {
	return config.DB.Save(&ip).Error
}

// DeleteAllowedIpByID deletes an AllowedIp by its ID
func DeleteAllowedIpByID(id string) error {
	if err := config.DB.Delete(&models.AllowedIp{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetAllowedIpsByService fetches AllowedIps by service (if specified)
func GetAllowedIpsByService(service string) ([]models.AllowedIp, error) {
	var ips []models.AllowedIp
	query := config.DB
	if service != "" {
		query = query.Where("service = ?", service)
	}
	if err := query.Find(&ips).Error; err != nil {
		return nil, errors.New("failed to fetch IPs")
	}
	return ips, nil
}
