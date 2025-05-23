package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddEmail(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Email  string ` binding:"required,email" json:"email"`
		Active bool   ` binding:"required" json:"active"`
	}

	var existingEmail models.SystemEmail

	if err := config.DB.First(&existingEmail).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "System Email already exists"})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sysEmail = models.SystemEmail{
		ID:     utils.GenerateUUID(),
		Email:  input.Email,
		Active: input.Active,
	}

	if err := config.DB.Create(&sysEmail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sysEmail)
}

func GetEmail(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var sysEmail models.SystemEmail

	if err := config.DB.First(&sysEmail).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": sysEmail})
}

func UpdateEmail(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Email  string ` binding:"required,email" json:"email"`
		Active bool   ` binding:"required" json:"active"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var sysEmail models.SystemEmail

	if err := config.DB.First(&sysEmail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	sysEmail.Active = input.Active
	sysEmail.Email = input.Email

	config.DB.Save(&sysEmail)

	c.JSON(http.StatusOK, sysEmail)
}
