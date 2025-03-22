package middlewares

import (
	"net/http"

	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/models"
	"github.com/Neimess/vpnbot_server/utils"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		telegramID, err := utils.GetTelegramId(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong TelegramID"})
			c.Abort()
			return
		}

		db := database.GetDB()
		var user models.User

		if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
