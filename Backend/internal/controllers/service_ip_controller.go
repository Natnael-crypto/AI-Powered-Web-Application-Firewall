package controllers

import (
	"net/http"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddAllowedIp(c *gin.Context) {
	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var input struct {
		Service string ` binding:"required,max=40" json:"service"`
		Ip      string ` binding:"required,ip" json:"ip"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowed_ip := models.AllowedIp{
		ID:      uuid.New().String(),
		Service: input.Service,
		Ip:      input.Ip,
	}
	if err := config.DB.Create(&allowed_ip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP"})
		return
	}
	c.JSON(http.StatusOK, allowed_ip)
}

func UpdateAllowedIp(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	var ip models.AllowedIp
	id := c.Param("id")
	if err := config.DB.Where("id = ?", id).First(&ip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
		return
	}

	var input struct {
		Ip string ` binding:"required,ip" json:"ip"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip.Ip = input.Ip
	config.DB.Save(&ip)

	c.JSON(http.StatusOK, ip)
}

func DeleteAllowedIp(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	id := c.Param("id")
	if err := config.DB.Delete(&models.AllowedIp{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "IP deleted"})
}

func GetAllowedIps(c *gin.Context) {

	if c.GetString("role") != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
		return
	}

	service := c.Query("service")
	var ips []models.AllowedIp
	query := config.DB
	if service != "" {
		query = query.Where("service = ?", service)
	}
	if err := query.Find(&ips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch IPs"})
		return
	}
	c.JSON(http.StatusOK, ips)
}
