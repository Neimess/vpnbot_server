package services

import (
	"errors"
	"log"
	"time"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/Neimess/vpnbot_server/wireguard"
)

func GetUserConfigs(telegramID int64) ([]models.Config, error) {
	db := database.GetDB()
	var user models.User
	var configs []models.Config

	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := db.Where("telegram_id = ?", telegramID).Find(&configs).Error; err != nil {
		return nil, err
	}

	return configs, nil
}

func CreateConfig(telegramID int64) (*models.Config, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	if !user.IsPaid {
		return nil, errors.New("user has not confirmed payment")
	}

	var profileCount int64
	db.Model(&models.Config{}).Where("telegram_id = ?", telegramID).Count(&profileCount)
	if profileCount >= 1 {
		return nil, errors.New("you have reached the limit of active profiles (1). Additional profiles will be available shortly")
	}

	privKey, pubKey, err := wireguard.GenerateKeys()
	if err != nil {
		return nil, err
	}

	clientIP, err := utils.GetFreeClientIP(db)
	if err != nil {
		return nil, err
	}
	configutaion, configPath, err := wireguard.GenerateClientConfig(db, clientIP, telegramID, privKey)
	if err != nil {
		return nil, err
	}

	if err := wireguard.AddClientToServer(telegramID, pubKey, clientIP); err != nil {
		return nil, err
	}

	config := models.Config{
		TelegramID: user.TelegramID,
		Config:     configutaion,
		ConfigPath: configPath,
		ClientIP:   clientIP,
		PublicKey:  pubKey,
		CreatedAt:  time.Now(),
	}

	if err := db.Create(&config).Error; err != nil {
		return nil, err
	}

	log.Println("Profile created for:", telegramID)
	return &config, nil
}

func DeleteUserConfig(configID uint, telegramID int64) error {
	db := database.GetDB()
	var user models.User
	var config models.Config

	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if err := db.Where("id = ? AND telegram_id = ?", configID, telegramID).First(&config).Error; err != nil {
		return errors.New("config not found")
	}

	if err := wireguard.RemoveClientServer(telegramID); err != nil {
		return err
	}
	if err := wireguard.RemoveClientData(config.ConfigPath); err != nil {
		return err
	}

	if err := db.Delete(&config).Error; err != nil {
		return err
	}

	return nil
}
