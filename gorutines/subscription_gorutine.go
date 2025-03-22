package gorutines

import (
	"log"
	"time"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/Neimess/vpnbot_server/wireguard"
)

func StartPaymentExpiryChecker() {

	ticker := time.NewTicker(time.Millisecond * 10000)
	go func() {
		for range ticker.C {
			log.Println("Starting daily payment expiry checker...")
			checkAndExpirePayments()
		}
	}()
}

func checkAndExpirePayments() {
	db := database.GetDB()
	var expiredUsers []models.User

	if err := db.Preload("Configs").Where("is_paid = ? AND expires_at <= ?", true, time.Now()).Find(&expiredUsers).Error; err != nil {
		log.Println("No expired users:", err)
		return
	}

	for _, user := range expiredUsers {
		log.Printf("Expiring VPN access for user: %d\n", user.TelegramID)

		for range user.Configs {
			err := wireguard.RemoveClientServer(user.TelegramID)
			if err != nil {
				log.Printf("Error removing user %d from wg0.conf: %v\n", user.TelegramID, err)
			}

		}

		if err := db.Where("telegram_id = ?", user.TelegramID).Delete(&models.Config{}).Error; err != nil {
			log.Printf("Error deleting configs from DB for user %d: %v\n", user.TelegramID, err)
		}

		user.IsPaid = false
		user.ExpiresAt = time.Time{}

		message := "Ваша подписка на VPN истекла. Чтобы продолжить использование, пожалуйста, оплатите продление."
		if err := utils.SendTelegramMessage(user.TelegramID, message); err != nil {
			log.Printf("Error sending message to user %d: %v", user.TelegramID, err)
		} else {
			log.Printf("Notification send to user: %d\n", user.TelegramID)
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user %d: %v\n", user.TelegramID, err)
		} else {
			log.Printf("VPN access expired and cleaned for user: %d\n", user.TelegramID)
		}
	}
}
