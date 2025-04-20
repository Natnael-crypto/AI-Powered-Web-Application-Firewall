package models

import (
	"time"
)

type SecurityHeader struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	ApplicationID string    `json:"application_id" gorm:"not null"`
	HeaderName    string    `json:"header_name" gorm:"unique;not null"`
	HeaderValue   string    `json:"header_value" gorm:"not null"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
