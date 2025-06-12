package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service     services.AttendanceService
	userService services.UserService
}

func NewAttendanceHandler(service services.AttendanceService, userService services.UserService) *AttendanceHandler {
	return &AttendanceHandler{
		service:     service,
		userService: userService,
	}
}

// GetAttendanceList godoc
// @Summary      Get attendance list
// @Description  Retrieves a list of attendance records for a specific user or all users if admin. Supports filtering by date range.
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        user_id   query     int  false  "User ID to filter attendance records"  default(0)
// @Param        start_date query     string  false  "Start date for filtering (YYYY-MM-DD)"
// @Param        end_date   query     string  false  "End date for filtering (YYYY-MM-DD)"
// @Success      200       {object}  map[string]interface{}  "List of attendance records"
// @Failure      400       {object}  map[string]string        "Invalid input"
// @Failure	     403       {object}  map[string]string        "Forbidden access"
// @Failure	     500       {object}  map[string]string        "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /attendances [get]
func (h *AttendanceHandler) GetAttendanceList(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in middleware"})
		return
	}
	currentUserIDUint, ok := currentUserID.(uint)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}
	// Check if the user is an admin or not
	isAdmin, err := h.userService.IsAdmin(currentUserIDUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check user role"})
		return
	}

	userIDParam := ctx.DefaultQuery("user_id", "0")
	startDate := ctx.DefaultQuery("start_date", "")
	endDate := ctx.DefaultQuery("end_date", "")
	if userIDParam == "0" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	// Check if the current user is not an admin and is trying to access another user's attendance
	if !isAdmin && uint(userID) != currentUserIDUint {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to access this user's attendance"})
		return
	}

	attendances, err := h.service.GetAttendanceList(uint(userID), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attendance list"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendances})
}

// RetrieveAttendance godoc
// @Summary      Retrieve attendance by ID
// @Description  Retrieves a specific attendance record by its ID. Admins can access any record, while users can only access their own.
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Attendance ID"
// @Success      200    {object}  map[string]interface{}  "Attendance record"
// @Failure      400    {object}  map[string]string        "Invalid ID"
// @Failure      404    {object}  map[string]string        "Attendance not found"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /attendances/{id} [get]
func (h *AttendanceHandler) RetrieveAttendance(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println("Attendance id : ", id)
	attendanceID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid attendance ID"})
		return
	}

	attendance, err := h.service.GetAttendanceByID(uint(attendanceID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "attendance not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendance})
}

// CheckIn godoc
// @Summary      Check in attendance
// @Description  Records a check-in for the current user. The user ID is obtained from the JWT token.
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Success      200    {object}  map[string]interface{}  "Attendance record"
// @Failure      400    {object}  map[string]string        "Invalid input"
// @Failure      401    {object}  map[string]string        "Unauthorized"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /attendances/check-in [post]
func (h *AttendanceHandler) CheckIn(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user access"})
		return
	}
	currentUserIDUint, ok := currentUserID.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id type"})
		return
	}

	userID := ctx.GetUint("user_id")
	requestID := ctx.GetString("request_id")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user_id", userID))
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_id", requestID))

	attendance, err := h.service.CheckIn(ctx, currentUserIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendance})
}

// CheckOut godoc
// @Summary      Check out attendance
// @Description  Records a check-out for a specific user. The user ID is provided as a query parameter.
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Success      200       {object}  map[string]interface{}  "Attendance record"
// @Failure      400       {object}  map[string]string        "Invalid input"
// @Failure      401       {object}  map[string]string        "Unauthorized"
// @Failure      403       {object}  map[string]string        "Forbidden access"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /attendances/check-out [post]
func (h *AttendanceHandler) CheckOut(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user access"})
		return
	}
	userID, ok := currentUserID.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	
	requestID := ctx.GetString("request_id")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user_id", userID))
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_id", requestID))

	attendance, err := h.service.CheckOut(ctx, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendance})
}
