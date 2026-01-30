package utils

import (
	"errors"
	"time"

	"example.com/event-booking/config"
	"github.com/golang-jwt/jwt/v5"
)

var cfg config.Config

func GenerateToken(email string, user_id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"user_id": user_id,
		"expired": time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(cfg.SecretKey))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(cfg.SecretKey), nil
	})

	if err != nil {
		return 0, errors.New("Failed to parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("Invalid token!")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid claim token!")
	}

	userId := int64(claims["user_id"].(float64))

	return userId, nil
}
