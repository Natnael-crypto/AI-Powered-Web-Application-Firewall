package models

import "time"

type Cert struct {
	CertID        string    `json:"cert_id" gorm:"primaryKey"`
	Cert          []byte    `json:"cert" gorm:"type:bytea"`
	Key           []byte    `json:"key" gorm:"type:bytea"`
	ApplicationID string    `json:"application_id" gorm:"unique;not null"`
	CreatedAt     time.Time `json:"created_at" `
	UpdatedAt     time.Time `json:"updated_at" `
}
