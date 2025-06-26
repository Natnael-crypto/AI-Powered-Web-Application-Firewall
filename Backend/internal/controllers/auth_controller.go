package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	resp, status := services.RegisterUserService(c)
	c.JSON(status, resp)
}

func LoginUser(c *gin.Context) {
	resp, status := services.LoginUserService(c)
	c.JSON(status, resp)
}

func UpdatePassword(c *gin.Context) {
	resp, status := services.UpdatePasswordService(c)
	c.JSON(status, resp)
}

func IsLoggedIN(c *gin.Context) {
	resp, status := services.IsLoggedInService(c)
	c.JSON(status, resp)
}
