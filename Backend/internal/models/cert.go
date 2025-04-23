package models

import "time"

type Cert struct {
	CertID        string    `json:"cert_id" gorm:"primaryKey"`
	Cert          []byte    `gorm:"type:bytea" gorm:"not null"  `
	Key           []byte    `gorm:"type:bytea" gorm:"not null" ` 
	ApplicationID string    `json:"application_id" gorm:"index" gorm:"unique;not null"`
	CreatedAt     time.Time `json:"created_at" `
	UpdatedAt     time.Time `json:"updated_at" `
}
