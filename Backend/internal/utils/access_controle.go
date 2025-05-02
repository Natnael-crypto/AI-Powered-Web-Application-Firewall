package utils

import (
	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAssignedApplicationIDs(c *gin.Context) []string {
	userID := c.GetString("user_id")
	if userID == "" {
		return []string{}
	}

	var mappings []models.UserToApplication
	if err := config.DB.Where("user_id = ?", userID).Find(&mappings).Error; err != nil {
		return []string{}
	}

	appIDs := make([]string, 0, len(mappings))
	for _, mapping := range mappings {
		appIDs = append(appIDs, mapping.ApplicationID)
	}

	return appIDs
}

func HasAccessToApplication(userAppIDs []string, targetAppID string) bool {
	for _, id := range userAppIDs {
		if id == targetAppID {
			return true
		}
	}
	return false
}
