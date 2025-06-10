package routes

import (
	"net/http"
	"time"

	"github.com/galiherlangga/go-attendance/app/handlers"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cache *redis.Client) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User modules
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router.GET("/health", func(ctx *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("login", userHandler.Login)
	}

	return router
}
