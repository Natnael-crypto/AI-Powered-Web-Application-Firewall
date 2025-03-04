package models

import "time"

type User struct {
	UserID               string    `json:"user_id" gorm:"primaryKey"`
	Username             string    `json:"username"`
	PasswordHash         string    `json:"password_hash"`
	Role                 string    `json:"role"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	LastLogin            time.Time `json:"last_login"`
	ProfileImageURL      string    `json:"profile_image_url"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
}
