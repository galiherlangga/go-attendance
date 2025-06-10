package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestComparePassword(t *testing.T) {
	t.Run("Comparing 2 password", func(t *testing.T) {
		password := "password"
		hashedPassword := "$2a$10$1YbpZQU8qK7IORyzyEayPeeTzHqzd2ULVawefaCdYybyvipqvoQHS" // Example hashed password

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			t.Errorf("Expected passwords to match, but they did not: %v", err)
		}
		assert.NoError(t, err, "Expected passwords to match")
	})
	
	t.Run("Comparing 2 different password", func(t *testing.T) {
		password := "differentpassword"
		hashedPassword := "$2a$10$EIX/5z1Zb7f8Q9e1j3k5uO0m6y5z5F4c5d3e7f8g9h0i1j2k3l4m5n" // Example hashed password

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		assert.Error(t, err, "Expected passwords to not match")
	})
}