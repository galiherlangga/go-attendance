package handlers

import (
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

func (h *AttendanceHandler) CheckIn(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in middleware"})
		return
	}
	currentUserIDUint, ok := currentUserID.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id type"})
		return
	}

	attendance, err := h.service.CheckIn(currentUserIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendance})
}

func (h *AttendanceHandler) CheckOut(ctx *gin.Context) {
	userIDParam := ctx.DefaultQuery("user_id", "0")
	if userIDParam == "0" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	attendance, err := h.service.CheckOut(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": attendance})
}
