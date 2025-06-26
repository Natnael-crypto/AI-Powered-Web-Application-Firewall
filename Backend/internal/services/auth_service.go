package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUserService(c *gin.Context) (gin.H, int) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	if repository.IsSuperAdminExists() {
		return gin.H{"error": "super admin exists"}, http.StatusUnauthorized
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return gin.H{"error": "failed to hash password"}, http.StatusInternalServerError
	}

	user := models.User{
		UserID:       utils.GenerateUUID(),
		Username:     input.Username,
		PasswordHash: hashedPassword,
		Role:         "super_admin",
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := repository.CreateUser(user); err != nil {
		return gin.H{"error": "username already exists"}, http.StatusConflict
	}

	return gin.H{"message": "user registered successfully"}, http.StatusCreated
}

func LoginUserService(c *gin.Context) (gin.H, int) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	user, err := repository.GetUserByUsername(input.Username)
	if err != nil || !utils.VerifyPassword(user.PasswordHash, input.Password) {
		return gin.H{"error": "invalid credentials"}, http.StatusUnauthorized
	}

	user.LastLogin = time.Now()
	if err := repository.UpdateUserLogin(user.UserID, user.LastLogin); err != nil {
		return gin.H{"error": "failed to update login time"}, http.StatusInternalServerError
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return gin.H{"error": "failed to generate token"}, http.StatusInternalServerError
	}

	return gin.H{"message": "login successful", "token": token}, http.StatusOK
}

func UpdatePasswordService(c *gin.Context) (gin.H, int) {
	var input models.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return gin.H{"error": err.Error()}, http.StatusBadRequest
	}

	user, err := repository.GetUserByUsername(input.Username)
	if err != nil || !utils.VerifyPassword(user.PasswordHash, input.OldPassword) {
		return gin.H{"error": "invalid credentials"}, http.StatusUnauthorized
	}

	newHash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return gin.H{"error": "failed to hash new password"}, http.StatusInternalServerError
	}

	if err := repository.UpdateUserPassword(user.UserID, newHash); err != nil {
		return gin.H{"error": "failed to update password"}, http.StatusInternalServerError
	}

	return gin.H{"message": "password updated successfully"}, http.StatusOK
}

func IsLoggedInService(c *gin.Context) (gin.H, int) {
	userID := c.GetString("user_id")
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return gin.H{"error": "Not logged in"}, http.StatusNotFound
	}
	user.PasswordHash = ""
	return gin.H{"user": user}, http.StatusOK
}
