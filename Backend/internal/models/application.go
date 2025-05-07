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
	// UserToApplication UserToApplication  `json:"user_to_application" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationID;references:ID"`
	// Cert              Cert               `json:"cert" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationID;references:ApplicationID"`
	// AppConf           AppConf            `json:"app_conf" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationID;references:ApplicationID"`
	// NotificationRules []NotificationRule `json:"notification_rules" gorm:"constraint:OnDelete:CASCADE;foreignKey:HostName;references:HostName"`
	// Requests          []Request          `json:"requests" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationName;references:HostName"`
	// Rules             []Rule             `json:"rules" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationID;references:ApplicationID"`
	// SecurityHeaders   []SecurityHeader   `json:"security_headers" gorm:"constraint:OnDelete:CASCADE;foreignKey:ApplicationID;references:ApplicationID"`
}

type UserToApplication struct {
	ID              string `json:"id" gorm:"primaryKey"`
	UserID          string `json:"user_id" gorm:"not null;index"`
	ApplicationID   string `json:"application_id" gorm:"not null"`
	ApplicationName string `json:"application_name"`
}
