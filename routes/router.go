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
	
	// Init repositories
	userRepo := repositories.NewUserRepository(db)
	payslipRepo := repositories.NewPayslipRepository(db)
	payrollPeriodRepo := repositories.NewPayrollPeriodRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	overtimeRepo := repositories.NewOvertimeRepository(db)
	reimbursementRepo := repositories.NewReimbursementRepository(db)
	
	// Init services
	userService := services.NewUserService(userRepo)
	payslipService := services.NewPayslipService(payslipRepo, attendanceRepo, overtimeRepo, reimbursementRepo, payrollPeriodRepo, userRepo)
	payrollPeriodService := services.NewPayrollPeriodService(payrollPeriodRepo, userRepo, payslipService, cache)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	overtimeService := services.NewOvertimeService(overtimeRepo, cache)
	reimbursementService := services.NewReimbursementService(reimbursementRepo, payrollPeriodRepo, cache)
	
	// Init handlers
	userHandler := handlers.NewUserHandler(userService)
	payslipHandler := handlers.NewPayslipHandler(payslipService, userService)
	payrollPeriodHandler := handlers.NewPayrollPeriodHandler(payrollPeriodService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService, userService)
	overtimeHandler := handlers.NewOvertimeHandler(overtimeService, userService)
	reimbursementHandler := handlers.NewReimbursementHandler(reimbursementService, userService)
	
	
	// Health check route
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
	
	// Payroll period routes
	payrollPeriodGroup := router.Group("/payroll-periods")
	payrollPeriodGroup.Use(middleware.IsAdminMiddleware(userRepo))
	{
		payrollPeriodGroup.GET("", payrollPeriodHandler.GetPayrollPeriodList)
		payrollPeriodGroup.GET("/:id", payrollPeriodHandler.GetPayrollPeriodByID)
		payrollPeriodGroup.POST("", payrollPeriodHandler.CreatePayrollPeriod)
		payrollPeriodGroup.PUT("/:id", payrollPeriodHandler.UpdatePayrollPeriod)
		payrollPeriodGroup.DELETE("/:id", payrollPeriodHandler.DeletePayrollPeriod)
		payrollPeriodGroup.POST("/:id/run-payroll", payrollPeriodHandler.RunPayrollPeriod)
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
	
	// Payslip routes
	payslipGroup := router.Group("/payslips")
	payslipGroup.Use(middleware.JWTAuthMiddleware())
	{
		payslipGroup.GET("/:period_id", payslipHandler.GetPayslipByUserAndPeriod)
	}
	payslipAdminGroup := router.Group("/payslips")
	payslipAdminGroup.Use(middleware.IsAdminMiddleware(userRepo))
	{
		payslipAdminGroup.GET("/summary/:period_id", payslipHandler.GetPayslipSummary)
	}

	return router
}
