package seeders

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/galiherlangga/go-attendance/app/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(db Database) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Println("Users already seeded, skipping...")
		return
	}

	fmt.Println("Seeding users...")
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating password: %v\n", err)
		return
	}
	adminUser := models.User{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: string(password),
		RoleID:   1, // Assuming role ID 1 is for admin
	}
	if err := db.Create(&adminUser).Error; err != nil {
		fmt.Printf("Error creating admin user: %v\n", err)
		return
	}
	fmt.Println("Admin user created successfully")

	// Loop to create 100 users
	for i := 1; i <= 100; i++ {
		user := models.User{
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: string(password), // Use the same password for simplicity
			RoleID:   2,                // Assuming role ID 2 is for regular users
		}
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Error creating user %d: %v\n", i, err)
			continue
		}
	}
}
