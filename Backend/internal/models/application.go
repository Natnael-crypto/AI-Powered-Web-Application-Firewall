package models

import "time"

type Application struct {
	ApplicationID   string    `json:"application_id" gorm:"primaryKey"`
	ApplicationName string    `json:"application_name"`
	Description     string    `json:"description"`
	HostName        string    `json:"hostname"`
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
