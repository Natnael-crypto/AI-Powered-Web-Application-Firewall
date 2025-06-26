package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddSecurityHeader(c *gin.Context) {
	resp, status := services.AddSecurityHeaderService(c)
	c.JSON(status, resp)
}

func GetSecurityHeaders(c *gin.Context) {
	resp, status := services.GetSecurityHeadersService(c)
	c.JSON(status, resp)
}

func GetSecurityHeadersAdmin(c *gin.Context) {
	resp, status := services.GetSecurityHeadersAdminService(c)
	c.JSON(status, resp)
}

func UpdateSecurityHeader(c *gin.Context) {
	resp, status := services.UpdateSecurityHeaderService(c)
	c.JSON(status, resp)
}

func DeleteSecurityHeader(c *gin.Context) {
	resp, status := services.DeleteSecurityHeaderService(c)
	c.JSON(status, resp)
}
