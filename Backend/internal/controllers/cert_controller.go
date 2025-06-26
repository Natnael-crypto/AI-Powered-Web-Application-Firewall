package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddCert(c *gin.Context) {
	resp, status := services.AddCertService(c)
	c.JSON(status, resp)
}

func GetCert(c *gin.Context) {
	resp, status := services.GetCertService(c)
	c.JSON(status, resp)
}

func GetCertAdmin(c *gin.Context) {
	resp, status := services.GetCertAdminService(c)
	c.JSON(status, resp)
}

func UpdateCert(c *gin.Context) {
	resp, status := services.UpdateCertService(c)
	c.JSON(status, resp)
}

func DeleteCert(c *gin.Context) {
	resp, status := services.DeleteCertService(c)
	c.JSON(status, resp)
}
