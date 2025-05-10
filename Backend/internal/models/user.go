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
	// UserToApplications   []UserToApplication `json:"user_to_applications" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	// Notifications        []Notification      `json:"notifications" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
	// NotificationConfig   NotificationConfig  `json:"notification_config" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID;references:UserID"`
}
