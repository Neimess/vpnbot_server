package database

import (
	"log"

	"errors"
	"time"

	"github.com/Neimess/vpnbot_server/config"
	"github.com/Neimess/vpnbot_server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	DB, err = gorm.Open(sqlite.Open(config.GlobalConfig.DATABASE_URI))
	if err != nil {
		log.Fatalf("Can't connect to db: %v", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Config{}, &models.Payment{})
	if err != nil {
		log.Fatalf("Error while migrating data: %v", err)
	}

	log.Println("Database successfully initialized")
}

func GetDB() *gorm.DB {
	return DB
}

func CreateAdmin(telegramID int64, name string) error {
	db := GetDB()

	var existingUser models.User

	err := db.Where("telegram_id = ?", telegramID).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAdmin := models.User{
				TelegramID: telegramID,
				Name:       name,
				IsAdmin:    true,
				IsPaid:     true,
				ExpiresAt:  time.Date(2050, 12, 31, 0, 0, 0, 0, time.UTC),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := db.Create(&newAdmin).Error; err != nil {
				return err
			}

			log.Printf("Admin with Telegram ID %d created successfully", telegramID)
			return nil
		}

		return err
	}

	// Если пользователь найден, обновляем его статус
	if err := db.Model(&existingUser).Updates(map[string]interface{}{
		"is_admin":   true,
		"is_paid":    true,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	log.Printf("User %d is now an admin", telegramID)
	return nil
}
