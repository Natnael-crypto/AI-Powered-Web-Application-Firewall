package models

type AllowedIp struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Service string `gorm:"not null" json:"service"`
	Ip      string `gorm:"not null" json:"ip"`
}
