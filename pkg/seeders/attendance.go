package seeders

import (
	"time"

	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

func SeedAttendances(db *gorm.DB) {
	// Load all users where role_id is 2 (regular users)
	var users []models.User
	if err := db.Where("role_id = ?", 2).Find(&users).Error; err != nil {
		panic("Failed to load users: " + err.Error())
	}
	checkinTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.Now().Location())
	checkoutTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, time.Now().Location())
	for _, user := range users {
		// Create attendance records for each user for the last 30 days
		for i := 0; i < 30; i++ {
			date := time.Now().AddDate(0, 0, -i)
			// Skip weekends (Saturday and Sunday)
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				continue
			}
			attendance := models.Attendance{
				UserID:    user.ID,
				Date:      date,
				CheckIn:   &checkinTime,
				CheckOut:  &checkoutTime,
				BaseModel: models.BaseModel{
					CreatedBy: &user.ID,
					UpdatedBy: &user.ID,
				},
			}
			if err := db.Create(&attendance).Error; err != nil {
				panic("Failed to create attendance record: " + err.Error())
			}
		}
	}
}
