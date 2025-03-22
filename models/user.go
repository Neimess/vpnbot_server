package models

import "time"

type User struct {
	TelegramID int64      `gorm:"primaryKey;autoIncrement:false" json:"telegram_id"`
	IsAdmin	   bool		  `json:"is_admin" gorm:"default:false"`
	Name       string     `gorm:"not null" json:"name"`
	IsPaid     bool       `json:"is_paid"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Configs    []Config   `gorm:"foreignKey:TelegramID;references:TelegramID" json:"configs"`
}
