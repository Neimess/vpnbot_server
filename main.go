package main

import (
	"log"
	"sync"

	"github.com/Neimess/vpnbot_server/config"
	"github.com/Neimess/vpnbot_server/database"
	"github.com/Neimess/vpnbot_server/gorutines"
	"github.com/Neimess/vpnbot_server/routes"
	"github.com/gin-gonic/gin"
)

func initialServer() {
	config.LoadConfig()

	database.InitDatabase()
	log.Println("Database succesfully connected")

	if err := database.CreateAdmin(int64(config.GlobalConfig.ADMIN_ID), config.GlobalConfig.ADMIN_NAME); err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}
}
func startServer(wg *sync.WaitGroup) {
	defer wg.Done()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies(nil)

	routes.RegisterUserRoutes(r)
	routes.RegisterAuthRoutes(r)
	routes.RegisterAdminRoutes(r)

	port := config.GlobalConfig.SERVER_PORT
	bind := "127.0.0.1:" + port
	log.Println("🚀 Server started at", bind)
	if err := r.Run(bind); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
func main() {
	initialServer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		gorutines.StartPaymentExpiryChecker()
	}()
	go startServer(&wg)
	wg.Wait()

}
