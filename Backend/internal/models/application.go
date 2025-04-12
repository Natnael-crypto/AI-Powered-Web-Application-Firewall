package models

import "time"

// Application represents the structure of an application in the system.
type Application struct {
	ApplicationID   string    `json:"application_id" gorm:"primaryKey"`
	ApplicationName string    `json:"application_name" `
	Description     string    `json:"description" `
	HostName        string    `json:"hostname" `
	IpAddress       string    `json:"ip_address" `
	Port            string    `json:"port" `
	Status          bool      `json:"status" `
	Tls             bool      `json:"tls" `
	CreatedAt       time.Time `json:"created_at" `
	UpdatedAt       time.Time `json:"updated_at" `
}

type UserToApplication struct {
	ID              string `json:"id" `
	UserID          string `json:"user_id" `
	ApplicationName string `json:"application_name" `
}
