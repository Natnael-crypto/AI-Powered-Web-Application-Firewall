package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAdminService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	var input struct {
		Username string `json:"username" binding:"required,min=4,max=12"`
		Password string `json:"password" binding:"required,min=8,max=25"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	if repository.AdminExists(input.Username) {
		return gin.H{"error": "username already exists"}, http.StatusConflict
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return gin.H{"error": "failed to hash password"}, http.StatusInternalServerError
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

	if err := repository.CreateAdmin(admin); err != nil {
		return gin.H{"error": "failed to save admin"}, http.StatusInternalServerError
	}

	return gin.H{"message": "admin added successfully"}, http.StatusCreated
}

func GetAdminService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	username := c.Param("username")
	admin, err := repository.GetAdminByUsername(username)
	if err != nil {
		return gin.H{"error": "admin not found"}, http.StatusNotFound
	}
	admin.PasswordHash = ""
	return gin.H{"admin": admin}, http.StatusOK
}

func GetAdminByIDService(c *gin.Context) (gin.H, int) {
	id := c.Param("user_id")

	if c.GetString("role") != "super_admin" && id != c.GetString("id") {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	admin, err := repository.GetAdminByID(id)
	if err != nil {
		return gin.H{"error": "admin not found"}, http.StatusNotFound
	}
	admin.PasswordHash = ""
	return gin.H{"admin": admin}, http.StatusOK
}

func GetAllAdminsService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	admins, err := repository.GetAllAdmins()
	if err != nil {
		return gin.H{"error": "unable to fetch admins"}, http.StatusInternalServerError
	}
	for i := range admins {
		admins[i].PasswordHash = ""
	}
	return gin.H{"admins": admins}, http.StatusOK
}

func DeleteAdminService(c *gin.Context) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	username := c.Param("username")
	if err := repository.DeleteAdminByUsername(username); err != nil {
		return gin.H{"error": "admin not found"}, http.StatusNotFound
	}
	return gin.H{"message": "admin deleted successfully"}, http.StatusOK
}

func UpdateAdminStatusService(c *gin.Context, status string) (gin.H, int) {
	if c.GetString("role") != "super_admin" {
		return gin.H{"error": "insufficient privileges"}, http.StatusForbidden
	}

	username := c.Param("username")
	if err := repository.UpdateAdminStatus(username, status); err != nil {
		return gin.H{"error": "failed to update admin status"}, http.StatusInternalServerError
	}
	return gin.H{"message": "admin status updated successfully"}, http.StatusOK
}
