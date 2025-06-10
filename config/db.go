package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/galiherlangga/go-attendance/pkg/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_USER", "gorm"),
		GetEnv("DB_PASS", "gorm"),
		GetEnv("DB_NAME", "gorm"),
		GetEnv("DB_PORT", "5432"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func InitTestDB() (*gorm.DB, error) {
	envPath, err := filepath.Abs(".env")
	if err != nil {
		log.Printf("Failed to resolve .env file path: %v\n", err)
		return nil, err
	}

	// Load environment variables from .env file
	err = LoadEnv(envPath)
	if err != nil {
		log.Printf("Failed to load .env file: %v\n", err)
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		GetEnv("TEST_DB_HOST", "localhost"),
		GetEnv("TEST_DB_USER", "gorm"),
		GetEnv("TEST_DB_PASS", "gorm"),
		GetEnv("TEST_DB_NAME", "gorm_test"),
		GetEnv("TEST_DB_PORT", "5432"),
	)
	log.Printf("Connecting to test database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migration once connected
	migrations.AutoMigrate(db)
	return db, nil
}
