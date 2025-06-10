package integrations

import (
	"testing"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/galiherlangga/go-attendance/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	t.Run("Successful Login", func(t *testing.T) {
		// Create a mock database
		mockDB := new(tests.MockDB)

		// Mock the Create method
		mockDB.On("Create", mock.Anything).Return(&gorm.DB{})

		// Create a test user
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		testUser := &models.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: string(password),
			RoleID:   2,
		}

		// Simulate saving the user
		result := mockDB.Create(testUser)

		// Assert that the Create method was called
		mockDB.AssertCalled(t, "Create", testUser)

		// Assert that the result is not nil
		assert.NotNil(t, result)

		// Simulate login logic
		// Mock finding the user by email
		mockDB.On("First", mock.Anything, "email = ?", testUser.Email).Return(&gorm.DB{})

		// Simulate password comparison
		err := bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte("password"))
		assert.NoError(t, err)
		
		// Simulate generating tokens
		accessToken, err := utils.GenerateAccessToken(testUser.ID)
		assert.NoError(t, err)
		refreshToken, err := utils.GenerateRefreshToken(testUser.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
	})
}
