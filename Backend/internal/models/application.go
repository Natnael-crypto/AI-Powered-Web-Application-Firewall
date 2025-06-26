package models

import "time"

type Application struct {
	ApplicationID   string    `json:"application_id" gorm:"primaryKey"`
	ApplicationName string    `json:"application_name"`
	Description     string    `json:"description"`
	HostName        string    `json:"host_name"`
	IpAddress       string    `json:"ip_address"`
	Port            string    `json:"port"`
	Status          bool      `json:"status"`
	Tls             bool      `json:"tls"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserToApplication struct {
	ID              string `json:"id" gorm:"primaryKey"`
	UserID          string `json:"user_id" gorm:"not null;index"`
	ApplicationID   string `json:"application_id" gorm:"not null"`
	ApplicationName string `json:"application_name"`
}

type ApplicationOptions struct {
	ApplicationID string `json:"application_id" `
	HostName      string `json:"hostname"`
}

type ApplicationInput struct {
	ApplicationName string `json:"application_name" binding:"required,max=20"`
	Description     string `json:"description" binding:"required,max=200"`
	HostName        string `json:"hostname" binding:"required,max=40"`
	IpAddress       string `json:"ip_address" binding:"required"`
	Port            string `json:"port" binding:"required,max=5"`
	Status          *bool  `json:"status" binding:"required"`
	Tls             *bool  `json:"tls" binding:"required"`
}
