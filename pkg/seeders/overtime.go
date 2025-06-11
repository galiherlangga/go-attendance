package seeders

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

func SeedOvertime(db *gorm.DB) {
	var users []models.User
	if err := db.Where("role_id = ?", 2).Find(&users).Error; err != nil {
		panic("Failed to load users: " + err.Error())
	}
	for _, user := range users {
		// Create overtime records for each user for the last 30 days
		for i := 0; i < 30; i++ {
			overtime := models.Overtime{
				UserID: user.ID,
				Date:   time.Now().AddDate(0, 0, -i),
				Hours:  gofakeit.IntRange(1, 3), // Example fixed hours for overtime
			}
			if err := db.Create(&overtime).Error; err != nil {
				panic("Failed to create overtime record: " + err.Error())
			}
		}
	}
}
