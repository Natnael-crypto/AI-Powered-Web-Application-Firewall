package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"github.com/gin-gonic/gin"
)

func AssignedApplications(c *gin.Context) []string {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if userID == "" || role == "" {
		return nil
	}

	var applicationIDs []string

	if role == "super_admin" {
		if err := config.DB.Model(&models.UserToApplication{}).
			Distinct("application_id").
			Pluck("application_id", &applicationIDs).Error; err != nil {
			return nil
		}
	} else {
		if err := config.DB.Model(&models.UserToApplication{}).
			Where("user_id = ?", userID).
			Pluck("application_id", &applicationIDs).Error; err != nil {
			return nil
		}
	}

	return applicationIDs
}
