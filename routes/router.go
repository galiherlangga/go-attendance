package routes

import (
	"net/http"
	"time"

	"github.com/galiherlangga/go-attendance/app/handlers"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/app/services"
	middleware "github.com/galiherlangga/go-attendance/pkg/middlewares"
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
	
	// Payroll period modules
	payrollPeriodRepo := repositories.NewPayrollPeriodRepository(db)
	payrollPeriodService := services.NewPayrollPeriodService(payrollPeriodRepo, cache)
	payrollPeriodHandler := handlers.NewPayrollPeriodHandler(payrollPeriodService)
	
	// Attendance modules
	attendanceRepo := repositories.NewAttendanceRepository(db)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, userService)
	
	// Overtime modules
	overtimeRepo := repositories.NewOvertimeRepository(db)
	overtimeService := services.NewOvertimeService(overtimeRepo, cache)
	overtimeHandler := handlers.NewOvertimeHandler(overtimeService, userService)
	
	// Reimbursement modules
	reimbursementRepo := repositories.NewReimbursementRepository(db)
	reimbursementService := services.NewReimbursementService(reimbursementRepo, payrollPeriodRepo, cache)
	reimbursementHandler := handlers.NewReimbursementHandler(reimbursementService, userService)
	
	// Auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("login", userHandler.Login)
	}
	
	// Payroll period routes
	payrollPeriodGroup := router.Group("/payroll-periods")
	payrollPeriodGroup.Use(middleware.IsAdminMiddleware(userRepo))
	{
		payrollPeriodGroup.GET("", payrollPeriodHandler.GetPayrollPeriodList)
		payrollPeriodGroup.GET("/:id", payrollPeriodHandler.GetPayrollPeriodByID)
		payrollPeriodGroup.POST("", payrollPeriodHandler.CreatePayrollPeriod)
		payrollPeriodGroup.PUT("/:id", payrollPeriodHandler.UpdatePayrollPeriod)
		payrollPeriodGroup.DELETE("/:id", payrollPeriodHandler.DeletePayrollPeriod)
	}
	
	// Attendance routes
	attendanceGroup := router.Group("/attendances")
	attendanceGroup.Use(middleware.JWTAuthMiddleware())
	{
		attendanceGroup.POST("/check-in", attendanceHandler.CheckIn)
		attendanceGroup.POST("/check-out", attendanceHandler.CheckOut)
		attendanceGroup.GET("", attendanceHandler.GetAttendanceList)
		attendanceGroup.GET("/:id", attendanceHandler.RetrieveAttendance)
	}
	
	// Overtime routes
	overtimeGroup := router.Group("/overtimes")
	overtimeGroup.Use(middleware.JWTAuthMiddleware())
	{
		overtimeGroup.GET("", overtimeHandler.GetOvertimeList)
		overtimeGroup.GET("/:id", overtimeHandler.GetOvertimeByID)
		overtimeGroup.POST("", overtimeHandler.CreateOvertime)
		overtimeGroup.PUT("/:id", overtimeHandler.UpdateOvertime)
		overtimeGroup.DELETE("/:id", overtimeHandler.DeleteOvertime)
	}
	
	// Reimbursement routes
	reimbursementGroup := router.Group("/reimbursements")
	reimbursementGroup.Use(middleware.JWTAuthMiddleware())
	{
		reimbursementGroup.GET("", reimbursementHandler.GetReimbursementList)
		reimbursementGroup.GET("/:id", reimbursementHandler.GetReimbursementByID)
		reimbursementGroup.POST("", reimbursementHandler.CreateReimbursement)
		reimbursementGroup.PUT("/:id", reimbursementHandler.UpdateReimbursement)
		reimbursementGroup.DELETE("/:id", reimbursementHandler.DeleteReimbursement)
	}

	return router
}
