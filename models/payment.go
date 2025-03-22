package models

import "time"

type Payment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TelegramID int64     `json:"telegram_id"`
	CreatedAt  time.Time `json:"created_at"`
}
