package routes

import (
	"github.com/Neimess/vpnbot_server/controllers"
	"github.com/Neimess/vpnbot_server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(router *gin.Engine) {
	admin := router.Group("/api/admin")
	admin.Use(middlewares.JWTMiddleware(), middlewares.AdminMiddleware())
	{
		admin.GET("/users", controllers.GetAllUsersHandler)
		admin.GET("/user", controllers.GetUserByIDHandler)
		admin.DELETE("/remove_user", controllers.AdminDeleteUserHandler)
		admin.PUT("/extend_any_subscription", controllers.AdminExtendSubscriptionHandler)
		admin.PUT("/create_config", controllers.AdminCreateConfigHandler)
		admin.GET("/payments", controllers.GetPaymentHistoryHandler)
	}
}
