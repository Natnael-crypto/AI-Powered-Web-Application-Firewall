package config

import (
	"backend/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateSuperAdminAccount() error {
	superAdminPassword := "super@admin123"
	hash, err := hashPassword(superAdminPassword)
	if err != nil {
		log.Fatalf("Super Admin seeding failed: %v", err)
	} else {
		log.Println("[+] Super Admin password:", superAdminPassword)
	}

	newSuperAdmin := models.User{
		UserID:          uuid.NewString(),
		Username:        "super admin",
		PasswordHash:    hash,
		Role:            "super_admin",
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       time.Unix(0, 0),
		ProfileImageURL: "",
	}

	var existingSuperAdmin models.User

	if err := DB.Where("role = ?", "super_admin").First(&existingSuperAdmin).Error; err == nil {
		return nil
	}

	if err := DB.Create(&newSuperAdmin).Error; err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}