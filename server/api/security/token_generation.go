package security

import (
	"crescendo-api/config/env"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := env.Load()
	if err != nil {
		log.Fatalf("Unable to load environmental variables file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("jwt secret not configured")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
