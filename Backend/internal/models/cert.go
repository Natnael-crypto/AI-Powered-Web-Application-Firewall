package models

import "time"

type Cert struct {
	CertID        string    `json:"cert_id" gorm:"primaryKey"`
	Cert          []byte    `gorm:"type:bytea"` // Store certificate as bytea for binary data in PostgreSQL
	Key           []byte    `gorm:"type:bytea"` // Store key as bytea for binary data in PostgreSQL
	ApplicationID string    `json:"application_id" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
