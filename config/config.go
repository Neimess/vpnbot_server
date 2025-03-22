package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DATABASE_URI      	string
	SERVER_IP         	string
	SERVER_PORT       	string
	WG_PORT           	string
	SERVER_PUBLIC_KEY 	string
	SERVER_URI			string
	JWT_SECRET        	string
	ADMIN_ID		  	int64
	ADMIN_NAME		  	string
	TELEGRAM_BOT_TOKEN	string
	ADMIN_SECRET		string
}

var GlobalConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using enviroments")
	}

	publicKey, err := GetWgPublicKey()
	if err != nil {
		log.Fatalf("WARNING: Could not retrieve WireGuard public key: %v", err)
	}

	adminID, err := strconv.ParseInt(getEnv("ADMIN_ID", "0"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid ADMIN_ID: %v", err)
	}

	GlobalConfig = Config{
		DATABASE_URI:      	getEnv("DATABASE_NAME", ""),
		SERVER_PORT:       	getEnv("SERVER_PORT", ""),
		SERVER_IP:         	getEnv("SERVER_IP", ""),
		WG_PORT:           	getEnv("WG_PORT", ""),
		JWT_SECRET:        	getEnv("JWT_SECRET", ""),
		ADMIN_ID:		   	adminID,
		ADMIN_NAME:		   	getEnv("ADMIN_NAME", "ADMIN"),
		TELEGRAM_BOT_TOKEN: getEnv("TELEGRAM_BOT_TOKEN", ""),
		ADMIN_SECRET:		getEnv("ADMIN_SECRET", ""),
		SERVER_PUBLIC_KEY: publicKey,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue == "" {
		log.Fatalf("ERROR: Required environment variable %s is missing", key)
	}
	return defaultValue
}
