package controllers

import (
	"net/http"

	"github.com/Neimess/vpnbot_server/config"
	"github.com/Neimess/vpnbot_server/services"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/gin-gonic/gin"
)

type RefreshTokenRequest struct {
	TelegramID  int64  `json:"telegram_id"`
	AdminSecret string `json:"admin_secret,omitempty"`
}

func RefreshTokenHandler(c *gin.Context) {
	var request RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	_, err := services.GetUser(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with current id doesn't exist"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func RefreshAdminTokenHandler(c *gin.Context) {
	var request RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	user, err := services.GetUser(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with current id doesn't exist"})
		return
	}

	if !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied. Not an admin"})
		return
	}

	if request.AdminSecret != config.GlobalConfig.ADMIN_SECRET {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin secret"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
