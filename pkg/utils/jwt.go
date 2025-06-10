package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAccessToken(userID uint) (string, error) {
	return generateJWT(userID, 15*time.Minute)
}

func GenerateRefreshToken(userID uint) (string, error) {
	return generateJWT(userID, 24*time.Hour)
}

func generateJWT(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
