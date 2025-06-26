package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddUserToApplication(c *gin.Context) {
	resp, status := services.AddUserToApplicationService(c)
	c.JSON(status, resp)
}

func UpdateUserToApplication(c *gin.Context) {
	resp, status := services.UpdateUserToApplicationService(c)
	c.JSON(status, resp)
}

func GetAllUserToApplications(c *gin.Context) {
	resp, status := services.GetAllUserToApplicationsService(c)
	c.JSON(status, resp)
}

func DeleteUserToApplication(c *gin.Context) {
	resp, status := services.DeleteUserToApplicationService(c)
	c.JSON(status, resp)
}
