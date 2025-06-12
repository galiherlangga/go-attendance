package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateAccessToken(userID uint) (string, error) {
	return generateJWT(userID, 15*time.Minute)
}

func GenerateRefreshToken(userID uint) (string, error) {
	return generateJWT(userID, 24*time.Hour)
}

func ParseJWT(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
	}
	return 0, fmt.Errorf("user_id not found in token claims")
}

func GetUserFromContext(ctx *gin.Context) (uint, error) {
	token, err := ctx.Cookie("access_token")
	if err != nil {
		// If cookie not found, check Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			return 0, fmt.Errorf("unauthorized – token missing")
		}
	}

	// Parse and validate the JWT token here
	userID, err := ParseJWT(token)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	if userID == 0 {
		return 0, fmt.Errorf("unauthorized – user ID not found")
	}
	return userID, nil
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
