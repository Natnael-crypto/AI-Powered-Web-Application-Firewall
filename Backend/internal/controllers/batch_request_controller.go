package controllers

import (
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func HandleBatchRequests(c *gin.Context) {
	resp, status := services.ProcessBatchRequestService(c)
	c.JSON(status, resp)
}
