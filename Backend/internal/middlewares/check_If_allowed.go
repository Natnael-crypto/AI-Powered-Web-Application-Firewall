package middleware

import (
	"net/http"
	"strings"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

func AllowlistMiddleware(c *gin.Context) {
	service := c.GetHeader("X-Service")
	ip := c.ClientIP()
	ipParts := strings.Split(ip, ":")

	if service == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-Service header"})
		c.Abort()
		return
	}

	var count int64
	config.DB.Model(&models.AllowedIp{}).Where("service = ? AND ip = ?", service, ipParts[0]).Count(&count)

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "IP not allowed for this service"})
		c.Abort()
		return
	}

	c.Next()
}