package routes

import (
	"github.com/Neimess/vpnbot_server/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	api := router.Group("/api/auth/")
	{
		api.POST("/refresh_token", controllers.RefreshTokenHandler)
		api.POST("/refresh_admin_token", controllers.RefreshAdminTokenHandler)
	}
}
