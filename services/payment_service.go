package services

import (
	"errors"
	"log"
	"time"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/wireguard"
)

func ConfirmPayment(telegramID int64) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Preload("Configs").Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	user.IsPaid = true
	user.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	if len(user.Configs) > 0 {
		for _, cfg := range user.Configs {
			err := wireguard.AddClientToServer(user.TelegramID, cfg.PublicKey, cfg.ClientIP)
			if err != nil {
				log.Printf("Failed to re-add config for user %d: %v", user.TelegramID, err)
			} else {
				log.Printf("Re-added config for user %d to wg0.conf", user.TelegramID)
			}
		}
	} else {
		log.Printf("No existing configs found for user %d; skipping wg0.conf restoration", user.TelegramID)
	}

	payment := models.Payment{
		TelegramID: telegramID,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(&payment).Error; err != nil {
		log.Printf("Error creating payment record: %v", err)
	}

	log.Println("Payment processed (create/extend) for:", telegramID)
	return &user, nil
}
