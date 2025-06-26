package repository

import (
	"backend/internal/config"
	"backend/internal/models"
)

func InsertBatchRequests(requests []models.Request) error {
	return config.DB.Create(&requests).Error
}

func GetApplicationHostMap() map[string]string {
	var applications []models.Application
	appMap := make(map[string]string)

	if err := config.DB.Find(&applications).Error; err != nil {
		return appMap
	}

	for _, app := range applications {
		appMap[app.HostName] = app.ApplicationID
	}
	return appMap
}
