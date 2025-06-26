package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddAdmin(c *gin.Context) {
	resp, status := services.AddAdminService(c)
	c.JSON(status, resp)
}

func GetAdmin(c *gin.Context) {
	resp, status := services.GetAdminService(c)
	c.JSON(status, resp)
}

func GetAdminByID(c *gin.Context) {
	resp, status := services.GetAdminByIDService(c)
	c.JSON(status, resp)
}

func GetAllAdmins(c *gin.Context) {
	resp, status := services.GetAllAdminsService(c)
	c.JSON(status, resp)
}

func DeleteAdmin(c *gin.Context) {
	resp, status := services.DeleteAdminService(c)
	c.JSON(status, resp)
}

func InactiveAdmin(c *gin.Context) {
	resp, status := services.UpdateAdminStatusService(c, "inactive")
	c.JSON(status, resp)
}

func ActiveAdmin(c *gin.Context) {
	resp, status := services.UpdateAdminStatusService(c, "active")
	c.JSON(status, resp)
}
