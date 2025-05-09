package models

type AllowedIp struct {
	ID      string `gorm:"primaryKey" json:"id"`
	Service string `gorm:"not null" json:"service"`
	Ip      string `gorm:"not null" json:"ip"`
}
