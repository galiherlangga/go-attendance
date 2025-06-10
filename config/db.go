package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_USER", "gorm"),
		GetEnv("DB_PASS", "gorm"),
		GetEnv("DB_NAME", "gorm"),
		GetEnv("DB_PORT", "9920"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
