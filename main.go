package main

import (
	"fmt"
	"log"

	"github.com/galiherlangga/go-attendance/config"
	"github.com/galiherlangga/go-attendance/pkg/migrations"
	"github.com/galiherlangga/go-attendance/pkg/seeders"
	"github.com/galiherlangga/go-attendance/routes"
)

// @title           Go Attendance API
// @version         1.0
// @description     This is the backend API for a Go-based attendance and payroll system.

// @contact.name   Galih Erlangga
// @contact.email  galiherlanggadev@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8010
// @BasePath  /

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadEnv(".env")
	appConfig := config.LoadAppConfig()
	
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
	migrations.AutoMigrate(db)
	seeders.Seed(db)
	
	config.InitRedis()
	
	r := routes.SetupRouter(db, config.RedisClient)
	
	addr := fmt.Sprintf("%s:%s", appConfig.Host, appConfig.Port)
	log.Printf("🚀 Server running at http://%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}