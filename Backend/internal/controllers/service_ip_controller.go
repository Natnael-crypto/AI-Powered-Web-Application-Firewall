package controllers

import (
	"net/http"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Add a new allowed IP
func AddAllowedIp(c *gin.Context) {
	var input models.AllowedIp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add IP"})
		return
	}
	c.JSON(http.StatusOK, input)
}

// Update an existing allowed IP
func UpdateAllowedIp(c *gin.Context) {
	var ip models.AllowedIp
	id := c.Param("id")
	if err := config.DB.First(&ip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "IP not found"})
		return
	}

	var input models.AllowedIp
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip.Service = input.Service
	ip.Ip = input.Ip
	config.DB.Save(&ip)

	c.JSON(http.StatusOK, ip)
}

// Delete an allowed IP
func DeleteAllowedIp(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.AllowedIp{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "IP deleted"})
}

// Get all allowed IPs or filter by service
func GetAllowedIps(c *gin.Context) {
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

// Check if IP is allowed for a given service
func CheckIfAllowed(c *gin.Context) {
	service := c.Query("service")
	ip := c.Query("ip")

	var count int64
	config.DB.Model(&models.AllowedIp{}).Where("service = ? AND ip = ?", service, ip).Count(&count)

	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"allowed": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"allowed": false})
	}
}
