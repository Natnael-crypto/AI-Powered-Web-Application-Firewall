package models

import "time"

type User struct {
	UserID          string    `json:"user_id" gorm:"primaryKey;unique;not null"`
	Username        string    `json:"username" gorm:"unique;not null"`
	PasswordHash    string    `json:"password_hash" gorm:"unique;not null"`
	Role            string    `json:"role" gorm:"not null"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	LastLogin       time.Time `json:"last_login"`
	ProfileImageURL string    `json:"profile_image_url"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordInput struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
