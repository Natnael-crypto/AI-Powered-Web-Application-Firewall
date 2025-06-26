package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

// GetConfig fetches the global configuration
func GetConfig(conf *models.Conf) error {
	if err := config.DB.First(conf).Error; err != nil {
		return err
	}
	return nil
}

// CreateConfig creates a new global configuration
func CreateConfig(conf models.Conf) error {
	if err := config.DB.Create(&conf).Error; err != nil {
		return err
	}
	return nil
}

// UpdateConfig updates an existing global configuration
func UpdateConfig(conf models.Conf) error {
	if err := config.DB.Save(&conf).Error; err != nil {
		return err
	}
	return nil
}

// GetAppConfig fetches the configuration for a specific application based on application ID
func GetAppConfig(applicationID string, appConf *models.AppConf) error {
	if err := config.DB.Where("application_id = ?", applicationID).First(appConf).Error; err != nil {
		return err
	}
	return nil
}

// // UpdateAppConfig updates an application's configuration
// func UpdateAppConfig(applicationID string, appConf models.AppConf) error {
// 	if err := config.DB.Where("application_id = ?", applicationID).Save(&appConf).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetAppConfigByID(applicationID string) (*models.AppConf, error) {
	var conf models.AppConf
	if err := config.DB.Where("application_id = ?", applicationID).First(&conf).Error; err != nil {
		return nil, err
	}
	return &conf, nil
}

func UpdateAppConfig(applicationID string, conf models.AppConf) error {
	if err := config.DB.Where("application_id = ?", applicationID).Save(&conf).Error; err != nil {
		return err
	}
	return nil
}
