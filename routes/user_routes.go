package routes

import (
	"github.com/Neimess/vpnbot_server/controllers"
	middleware "github.com/Neimess/vpnbot_server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/create_user", controllers.CreateUserHandler)
		protected := api.Group("/")
		protected.Use(middleware.JWTMiddleware())
		{
			protected.GET("/get_user", controllers.GetUserHandler)
			protected.PUT("/confirm_payment", controllers.ConfirmPaymentHandler)
			protected.POST("/create_config", controllers.CreateUserConfigHandler)
			protected.GET("/get_configs", controllers.GetUserConfigsHandler)
			protected.DELETE("/del_config", controllers.DeleteUserConfigHandler)
		}
	}
}
