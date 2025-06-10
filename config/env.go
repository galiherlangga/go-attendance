package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(filename string) error {
	err := godotenv.Load(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func GetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}