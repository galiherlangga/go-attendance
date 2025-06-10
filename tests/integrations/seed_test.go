package integrations

import (
	"testing"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/config"
	"github.com/galiherlangga/go-attendance/pkg/seeders"
	"github.com/stretchr/testify/assert"
)

func TestSeed(t *testing.T) {
	t.Run("Successful Seeding", func(t *testing.T) {
		// Initialize the test database
		db, err := config.InitTestDB()
		assert.NoError(t, err)
		defer func() {
			sqlDB, _ := db.DB()
			sqlDB.Close()
		}()

		// Run the seeding function
		seeders.Seed(db)

		// Verify the admin user was created
		var adminUser models.User
		err = db.Where("email = ?", "admin@example.com").First(&adminUser).Error
		assert.NoError(t, err)
		assert.Equal(t, "Admin", adminUser.Name)
		assert.Equal(t, "admin@example.com", adminUser.Email)

		// Verify 100 regular users were created
		var userCount int64
		err = db.Model(&models.User{}).Where("role_id = ?", 2).Count(&userCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(100), userCount)
	})
}
