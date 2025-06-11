package migrations

import (
	"log"

	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.PayrollPeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payslip{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")
}