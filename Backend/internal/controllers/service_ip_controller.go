package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddAllowedIp(c *gin.Context) {
	resp, status := services.AddAllowedIpService(c)
	c.JSON(status, resp)
}

func UpdateAllowedIp(c *gin.Context) {
	resp, status := services.UpdateAllowedIpService(c)
	c.JSON(status, resp)
}

func DeleteAllowedIp(c *gin.Context) {
	resp, status := services.DeleteAllowedIpService(c)
	c.JSON(status, resp)
}

func GetAllowedIps(c *gin.Context) {
	resp, status := services.GetAllowedIpsService(c)
	c.JSON(status, resp)
}
