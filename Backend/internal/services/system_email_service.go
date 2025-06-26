package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddEmailService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		Email  string `binding:"required,email" json:"email"`
		Active bool   `binding:"required" json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	_, err := repository.GetSystemEmail()
	if err == nil {
		return gin.H{"error": "System Email already exists"}, http.StatusConflict
	}

	sysEmail := models.SystemEmail{
		ID:     utils.GenerateUUID(),
		Email:  input.Email,
		Active: input.Active,
	}

	if err := repository.CreateSystemEmail(sysEmail); err != nil {
		return gin.H{"error": err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"email": sysEmail}, http.StatusOK
}

func GetEmailService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	sysEmail, err := repository.GetSystemEmail()
	if err != nil {
		return gin.H{"error": err.Error()}, http.StatusNotFound
	}

	return gin.H{"email": sysEmail}, http.StatusOK
}

func UpdateEmailService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		Email  string `binding:"required,email" json:"email"`
		Active bool   `binding:"required" json:"active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	sysEmail, err := repository.GetSystemEmail()
	if err != nil {
		return gin.H{"error": err.Error()}, http.StatusInternalServerError
	}

	sysEmail.Email = input.Email
	sysEmail.Active = input.Active

	if err := repository.UpdateSystemEmail(sysEmail); err != nil {
		return gin.H{"error": err.Error()}, http.StatusInternalServerError
	}

	return gin.H{"email": sysEmail}, http.StatusOK
}
