package seeders

import (
	"fmt"

	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	var count int64
	db.Model(&models.Role{}).Count(&count)
	if count > 0 {
		fmt.Println("Roles already seeded, skipping...")
		return
	}
	roles := []string{"admin", "user"}
	
	fmt.Println("Seeding roles...")
	var roleModels []models.Role
	for _, roleName := range roles {
		roleModels = append(roleModels, models.Role{Name: roleName})
	}

	db.Create(&roleModels)
}