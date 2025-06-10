package units

import (
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	duration := 15 * 60 * 1000 // 15 minutes in milliseconds
	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Duration(duration) * time.Millisecond).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotEmpty(t, tokenString, "Token should not be empty")
	log.Println("Generated JWT Token:", tokenString)
}
