package repository

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"time"
)

func ApplicationExistsByName(name string) bool {
	var app models.Application
	return config.DB.Where("application_name = ?", name).First(&app).Error == nil
}

func ApplicationExistsByHost(host string) bool {
	var app models.Application
	return config.DB.Where("hostname = ?", host).First(&app).Error == nil
}

func CreateApplication(app models.Application) error {
	return config.DB.Create(&app).Error
}

func CreateAppConfig(conf models.AppConf) error {
	return config.DB.Create(&conf).Error
}

func CreateUserToApp(uta models.UserToApplication) error {
	return config.DB.Create(&uta).Error
}

func CreateEmptyCertificate(appID string) {
	cert := models.Cert{
		CertID:        utils.GenerateUUID(),
		ApplicationID: appID,
		Cert:          []byte{},
		Key:           []byte{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	config.DB.Create(&cert)
}

func GetApplicationByID(id string) (models.Application, error) {
	var app models.Application
	err := config.DB.Where("application_id = ?", id).First(&app).Error
	return app, err
}

func FetchAllApplicationsWithConfig() ([]map[string]interface{}, error) {
	var apps []models.Application
	if err := config.DB.Find(&apps).Error; err != nil {
		return nil, err
	}
	return mapAppsWithConfigs(apps)
}

func FetchApplicationsByIDsWithConfig(ids []string) ([]map[string]interface{}, error) {
	var apps []models.Application
	if err := config.DB.Where("application_id IN ?", ids).Find(&apps).Error; err != nil {
		return nil, err
	}
	return mapAppsWithConfigs(apps)
}

func mapAppsWithConfigs(apps []models.Application) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for _, app := range apps {
		appMap := map[string]interface{}{
			"application_id":   app.ApplicationID,
			"application_name": app.ApplicationName,
			"description":      app.Description,
			"hostname":         app.HostName,
			"ip_address":       app.IpAddress,
			"port":             app.Port,
			"status":           app.Status,
			"tls":              app.Tls,
			"created_at":       app.CreatedAt,
			"updated_at":       app.UpdatedAt,
		}

		var appConf models.AppConf
		if err := config.DB.Where("application_id = ?", app.ApplicationID).First(&appConf).Error; err == nil {
			appMap["config"] = appConf
		} else {
			appMap["config"] = nil
		}
		result = append(result, appMap)
	}
	return result, nil
}

func UpdateApplication(appID string, input models.ApplicationInput) error {
	return config.DB.Model(&models.Application{}).Where("application_id = ?", appID).Updates(map[string]interface{}{
		"application_name": input.ApplicationName,
		"description":      input.Description,
		"host_name":        input.HostName,
		"ip_address":       input.IpAddress,
		"port":             input.Port,
		"status":           *input.Status,
		"tls":              *input.Tls,
		"updated_at":       time.Now(),
	}).Error
}

func DeleteApplicationAndLinks(appID string) error {
	if err := config.DB.Where("application_id = ?", appID).Delete(&models.Application{}).Error; err != nil {
		return err
	}
	return config.DB.Where("application_id = ?", appID).Delete(&models.UserToApplication{}).Error
}
