package models

import (
	"time"
)

type Config struct {
	ID         		uint      	`gorm:"primaryKey" json:"id"`
	TelegramID     	int64      	`gorm:"index" json:"telegram_id"`
	ConfigPath 		string    	`json:"config_path"`
	Config	   		string	  	`json:"config"`
	ClientIP    	string    	`gorm:"uniqueIndex;not null" json:"client_ip"`
	PublicKey  		string    	`json:"public_key"`
	CreatedAt  		time.Time 	`json:"created_at"`
}	