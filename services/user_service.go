package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/utils"
	"gorm.io/gorm"
)

func CreateUser(telegramID int64, name string) (*models.User, string, error) {
	db := database.GetDB()
	var existingUser models.User

	if err := db.Where("telegram_id = ?", telegramID).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User not found, proceeding to create new user")
		} else {
			log.Printf("Database error while checking user existence: %v", err)
			return nil, "", fmt.Errorf("database error: %v", err)
		}
	} else {
		log.Printf("User with Telegram ID %d already exists", telegramID)
		return nil, "", fmt.Errorf("user with Telegram ID %d already exists", telegramID)
	}

	user := &models.User{
		TelegramID: telegramID,
		Name:       name,
		IsPaid:     false,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, "", err
	}
	accessToken, err := utils.GenerateAccessToken(telegramID)
	if err != nil {
		return nil, "", err
	}
	log.Println("Create user", telegramID)
	return user, accessToken, nil
}

func GetUser(telegramID int64) (*models.User, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Preload("Configs").First(&user, "telegram_id = ?", telegramID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(telegramID int64) error {
	db := database.GetDB()
	var user models.User

	if result := db.Where("telegram_id = ?", telegramID).Delete(&user).RowsAffected; result == 0 {
		return errors.New("user not found")
	}

	return nil
}
