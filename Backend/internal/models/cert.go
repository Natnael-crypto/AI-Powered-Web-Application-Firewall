package models

import "time"

type Cert struct {
	CertID        string    `json:"cert_id" gorm:"primaryKey"`
	Cert          []byte    `gorm:"type:bytea" gorm:"not null"  ` // Store certificate as bytea for binary data in PostgreSQL
	Key           []byte    `gorm:"type:bytea" gorm:"not null" `  // Store key as bytea for binary data in PostgreSQL
	ApplicationID string    `json:"application_id" gorm:"index" gorm:"unique;not null"`
	CreatedAt     time.Time `json:"created_at" `
	UpdatedAt     time.Time `json:"updated_at" `
}
