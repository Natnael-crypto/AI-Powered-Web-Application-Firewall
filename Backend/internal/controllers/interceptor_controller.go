package controllers

import (
	"backend/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckForUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"change": config.Change})
	config.Change = false
}
