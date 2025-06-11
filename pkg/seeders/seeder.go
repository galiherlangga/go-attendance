package seeders

import "gorm.io/gorm"

type Database interface {
	Model(value interface{}) *gorm.DB
	Count(count *int64) *gorm.DB
	Create(value interface{}) *gorm.DB
}

func Seed(db *gorm.DB) {
	SeedRoles(db)
	SeedUsers(db)
	SeedAttendances(db)
	SeedOvertime(db)
}