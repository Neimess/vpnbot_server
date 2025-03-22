package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Neimess/vpnbot_server/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTelegramId(c *gin.Context) (int64, error) {
	telegramIDInterface, exists := c.Get("telegram_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0, errors.New("payload doesn't exists")
	}
	telegramID, ok := telegramIDInterface.(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid telegram_id format"})
		return 0, errors.New("error while converting float to int")
	}
	return telegramID, nil
}

func GetFreeClientIP(db *gorm.DB) (string, error) {
	var configs []models.Config
	if err := db.Select("client_ip").Find(&configs).Error; err != nil {
		return "", fmt.Errorf("failed to get existing IPs: %w", err)
	}

	used := make(map[string]bool)
	for _, cfg := range configs {
		used[cfg.ClientIP] = true
	}

	for i := 2; i < 255; i++ {
		ip := fmt.Sprintf("10.0.0.%d/32", i)
		if !used[ip] {
			return ip, nil
		}
	}

	return "", fmt.Errorf("no available IP addresses")
}
