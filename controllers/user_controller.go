package controllers

import (
	// "fmt"
	"net/http"

	"github.com/Neimess/vpnbot_server/services"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Name       string `json:"name"`
}

type ExtendSubscriptionRequest struct {
	Days int `json:"days"`
}

func CreateUserHandler(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "unavailable JSON format"})
		return
	}

	user, access_token, err := services.CreateUser(request.TelegramID, request.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":         user,
		"access_token": access_token,
	})
}

func GetUserHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
		return
	}
	user, err := services.GetUser(telegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUserHandler(c *gin.Context) {
	telegramID, err := utils.GetTelegramId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
		return
	}

	err = services.DeleteUser(telegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK,
		gin.H{
			"message":     "User was succesfully deleted",
			"telegram_id": telegramID,
		})
}
