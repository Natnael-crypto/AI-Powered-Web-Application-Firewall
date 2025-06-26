package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddAllowedIpService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		Service string `binding:"required,max=40" json:"service"`
		Ip      string `binding:"required,ip" json:"ip"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	allowedIp := models.AllowedIp{
		ID:      uuid.New().String(),
		Service: input.Service,
		Ip:      input.Ip,
	}
	if err := repository.CreateAllowedIp(allowedIp); err != nil {
		return gin.H{"error": "Failed to add IP"}, http.StatusInternalServerError
	}

	if input.Service == "I" {
		config.Change = true
	}
	return gin.H{"allowed_ip": allowedIp}, http.StatusOK
}

func UpdateAllowedIpService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	id := c.Param("id")
	ip, err := repository.GetAllowedIpByID(id)
	if err != nil {
		return gin.H{"error": "IP not found"}, http.StatusNotFound
	}

	var input struct {
		Ip string `binding:"required,ip" json:"ip"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	ip.Ip = input.Ip
	if err := repository.UpdateAllowedIp(ip); err != nil {
		return gin.H{"error": "Failed to update IP"}, http.StatusInternalServerError
	}

	if ip.Service == "I" {
		config.Change = true
	}
	return gin.H{"allowed_ip": ip}, http.StatusOK
}

func DeleteAllowedIpService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	id := c.Param("id")
	if err := repository.DeleteAllowedIpByID(id); err != nil {
		return gin.H{"error": "Failed to delete IP"}, http.StatusInternalServerError
	}

	config.Change = true
	return gin.H{"message": "IP deleted successfully"}, http.StatusOK
}

func GetAllowedIpsService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	service := c.Query("service")
	ips, err := repository.GetAllowedIpsByService(service)
	if err != nil {
		return gin.H{"error": "Failed to fetch IPs"}, http.StatusInternalServerError
	}

	return gin.H{"allowed_ips": ips}, http.StatusOK
}
