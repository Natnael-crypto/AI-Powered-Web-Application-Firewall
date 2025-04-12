package models

import (
	"time"
)

type SecurityHeader struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	ApplicationID string    `json:"application_id"`
	HeaderName    string    `json:"header_name"`
	HeaderValue   string    `json:"header_value"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
