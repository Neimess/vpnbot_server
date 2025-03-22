package controllers

import (
	"net/http"
	"time"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/services"
	"github.com/gin-gonic/gin"
)

type AdminRequestID struct {
	TelegramID int64   `json:"telegram_id"`
	Amount     float64 `json:"amount"`
}

func GetAllUsersHandler(c *gin.Context) {
	db := database.GetDB()
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByIDHandler(c *gin.Context) {
	var request AdminRequestID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("telegram_id = ?", request.TelegramID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func AdminDeleteUserHandler(c *gin.Context) {
	var request AdminRequestID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("telegram_id = ?", request.TelegramID).Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func AdminExtendSubscriptionHandler(c *gin.Context) {
	var request AdminRequestID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	db := database.GetDB()
	var user models.User

	if err := db.Where("telegram_id = ?", request.TelegramID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	newDate := user.ExpiresAt.AddDate(0, 1, 0)
	if user.ExpiresAt.Before(time.Now()) {
		newDate = time.Now().AddDate(0, 1, 0)
	}

	if err := db.Model(&user).Update("expires_at", newDate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extend subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func AdminCreateConfigHandler(c *gin.Context) {
	var request AdminRequestID

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	_, err := services.GetUser(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User with current id doesn't exists"})
		return
	}

	_, err = services.ConfirmPayment(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	config, err := services.CreateConfig(request.TelegramID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Config created successfully, user is marked as paid",
		"telegram_id": request.TelegramID,
		"is_paid":     true,
		"config":      config,
	})
}

func GetPaymentHistoryHandler(c *gin.Context) {
	db := database.GetDB()
	var payments []models.Payment

	if err := db.Order("created_at desc").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment history"})
		return
	}

	c.JSON(http.StatusOK, payments)
}
