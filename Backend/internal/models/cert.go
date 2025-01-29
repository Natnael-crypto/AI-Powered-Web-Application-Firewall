package models

import "time"

type Cert struct {
	CertID        string    `json:"cert_id" gorm:"primaryKey"`
	CertPath      string    `json:"cert_path"`
	KeyPath       string    `json:"key_path"`
	ApplicationID string    `json:"application_id" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
