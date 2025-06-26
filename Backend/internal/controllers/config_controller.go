package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	resp, status := services.GetConfigService()
	c.JSON(status, resp)
}

func GetConfigAdmin(c *gin.Context) {
	resp, status := services.GetConfigAdminService()
	c.JSON(status, resp)
}

func GetAppConfig(c *gin.Context) {
	resp, status := services.GetAppConfigService(c)
	c.JSON(status, resp)
}

func CreateConfig(c *gin.Context) {
	resp, status := services.CreateConfigService(c)
	c.JSON(status, resp)
}

func UpdateListeningPort(c *gin.Context) {
	resp, status := services.UpdateListeningPortService(c)
	c.JSON(status, resp)
}

func UpdateRateLimit(c *gin.Context) {
	resp, status := services.UpdateRateLimitService(c)
	c.JSON(status, resp)
}

func UpdateTls(c *gin.Context) {
	resp, status := services.UpdateTlsService(c)
	c.JSON(status, resp)
}

func UpdateDetectBot(c *gin.Context) {
	resp, status := services.UpdateDetectBotService(c)
	c.JSON(status, resp)
}

func UpdateRemoteLogServer(c *gin.Context) {
	resp, status := services.UpdateRemoteLogServerService(c)
	c.JSON(status, resp)
}

func UpdateMaxPostDataSize(c *gin.Context) {
	resp, status := services.UpdateMaxPostDataSizeService(c)
	c.JSON(status, resp)
}
