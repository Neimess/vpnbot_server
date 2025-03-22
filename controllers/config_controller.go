package controllers

import (
	"net/http"
	"strconv"

	"github.com/Neimess/vpnbot_server/services"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/gin-gonic/gin"
)

func GetUserConfigsHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
		return
	}

	configs, err := services.GetUserConfigs(telegramID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func DeleteUserConfigHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
		return
	}

	configID, err := strconv.ParseUint(c.Param("config_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config_id"})
		return
	}

	err = services.DeleteUserConfig(uint(configID), telegramID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "config deleted successfully"})
}

func CreateUserConfigHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
		return
	}

	user, err := services.GetUser(telegramID)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if !user.IsPaid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Subscription not active. Please pay first."})
		return
	}

	config, err := services.CreateConfig(telegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to create config: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config created successfully",
		"config":  config,
	})
}
