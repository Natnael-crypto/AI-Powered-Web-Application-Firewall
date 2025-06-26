package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEmail(c *gin.Context) {
	resp, status := services.AddEmailService(c)
	c.JSON(status, resp)
}

func GetEmail(c *gin.Context) {
	resp, status := services.GetEmailService(c)
	c.JSON(status, resp)
}

func UpdateEmail(c *gin.Context) {
	resp, status := services.UpdateEmailService(c)
	c.JSON(status, resp)
}
