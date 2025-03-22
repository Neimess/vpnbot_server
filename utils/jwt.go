package utils

import (
	"fmt"
	"time"

	"github.com/Neimess/vpnbot_server/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(telegramID int64) (string, error) {
	accessClaims := jwt.MapClaims{
		"telegram_id": telegramID,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(config.GlobalConfig.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return accessString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.GlobalConfig.JWT_SECRET), nil
	})
}

func ExtractClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
