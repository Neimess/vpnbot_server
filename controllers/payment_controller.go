package controllers

import (
	"net/http"

	"github.com/Neimess/vpnbot_server/services"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/gin-gonic/gin"
)

func ConfirmPaymentHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	user, err := services.ConfirmPayment(telegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment processed", "user": user})
}
