package config

import (
	"log"
	"os"
)

func GetWgPublicKey() (string, error) {
	data, err := os.ReadFile("/configs/server/publickey-server")
	if err != nil {
		return "", err
	}
	log.Println(data)
	return string(data), nil
}
