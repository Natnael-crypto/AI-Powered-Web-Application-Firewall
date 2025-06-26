package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddApplication(c *gin.Context) {
	resp, status := services.AddApplicationService(c)
	c.JSON(status, resp)
}

func GetApplication(c *gin.Context) {
	resp, status := services.GetApplicationService(c)
	c.JSON(status, resp)
}

func GetAllApplications(c *gin.Context) {
	resp, status := services.GetAllApplicationsService(c)
	c.JSON(status, resp)
}

func GetAllApplicationsAdmin(c *gin.Context) {
	resp, status := services.GetAdminApplicationsService(c)
	c.JSON(status, resp)
}

func UpdateApplication(c *gin.Context) {
	resp, status := services.UpdateApplicationService(c)
	c.JSON(status, resp)
}

func DeleteApplication(c *gin.Context) {
	resp, status := services.DeleteApplicationService(c)
	c.JSON(status, resp)
}
