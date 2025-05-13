package controllers

import (
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func AddAdmin(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Username string `json:"username" binding:"required,min=4,max=12"`
		Password string `json:"password" binding:"required,min=8,max=25"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingAdmin models.User
	if err := config.DB.Where("username = ?", input.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	admin := models.User{
		UserID:       utils.GenerateUUID(),
		Username:     input.Username,
		PasswordHash: hashedPassword,
		Role:         "admin",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save admin"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "admin added successfully"})
}

func GetAdmin(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	username := c.Param("username")

	var admin models.User
	if err := config.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	admin.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"admin": admin})
}

func GetAdminByID(c *gin.Context) {
	id := c.Param("user_id")

	if c.GetString("role") == "super_admin" {
	} else {
		if id != c.GetString("id") {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
			return
		}
	}

	var admin models.User
	if err := config.DB.Where("user_id = ?", id).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	admin.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"admin": admin})
}

func GetAllAdmins(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var admins []models.User
	if err := config.DB.Where("role = ?", "admin").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch admins"})
		return
	}

	for i := range admins {
		admins[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, gin.H{"admins": admins})
}

func DeleteAdmin(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	username := c.Param("username")

	if err := config.DB.Where("username = ?", username).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin deleted successfully"})
}

func InactiveAdmin(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	username := c.Param("username")

	if err := config.DB.Model(&models.User{}).Where("username = ?", username).Update("status", "inactive").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin status updated successfully"})
}

func ActiveAdmin(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	username := c.Param("username")

	if err := config.DB.Model(&models.User{}).Where("username = ?", username).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin status updated successfully"})
}
