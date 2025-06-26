package repository

import (
	"backend/internal/config"
	"backend/internal/models"
	"time"
)

func IsSuperAdminExists() bool {
	var user models.User
	err := config.DB.Where("role = ?", "super_admin").First(&user).Error
	return err == nil
}

func CreateUser(user models.User) error {
	return config.DB.Create(&user).Error
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// func GetUserByID(userID string) (models.User, error) {
// 	var user models.User
// 	err := config.DB.Where("user_id = ?", userID).First(&user).Error
// 	return user, err
// }

func UpdateUserLogin(userID string, lastLogin time.Time) error {
	return config.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"last_login": lastLogin,
		"updated_at": time.Now(),
	}).Error
}

func UpdateUserPassword(userID, newHash string) error {
	return config.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"password_hash": newHash,
		"updated_at":    time.Now(),
	}).Error
}
