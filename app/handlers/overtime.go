package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/gin-gonic/gin"
)

type OvertimeHandler struct {
	service     services.OvertimeService
	userService services.UserService
}

func NewOvertimeHandler(service services.OvertimeService, userService services.UserService) *OvertimeHandler {
	return &OvertimeHandler{
		service:     service,
		userService: userService,
	}
}

func (h *OvertimeHandler) GetOvertimeList(ctx *gin.Context) {
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

	pagination := utils.GetPagination(ctx)

	overtimeList, total, err := h.service.GetOvertimeList(uint(userID), pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve overtime list"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":       overtimeList,
		"total":      total,
		"page":       pagination.Page,
		"limit":      pagination.Limit,
		"totalPages": int(total / int64(pagination.Limit)),
	})

}

func (h *OvertimeHandler) GetOvertimeByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	overtime, err := h.service.GetOvertimeByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Overtime not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": overtime})
}

func (h *OvertimeHandler) CreateOvertime(ctx *gin.Context) {
	var overtimeReq models.OvertimeRequest
	if err := ctx.ShouldBindJSON(&overtimeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

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
	overtime := models.Overtime{
		UserID: currentUserIDUint,
		Date:   overtimeReq.Date,
		Hours:  overtimeReq.Hours,
		Note:   overtimeReq.Note,
	}

	createdOvertime, err := h.service.SubmitOvertime(&overtime)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create overtime", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdOvertime})
}

func (h *OvertimeHandler) UpdateOvertime(ctx *gin.Context) {
	var overtimeReq models.OvertimeRequest
	if err := ctx.ShouldBindJSON(&overtimeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

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
	overtime := models.Overtime{
		UserID: currentUserIDUint,
		Date:   overtimeReq.Date,
		Hours:  overtimeReq.Hours,
		Note:   overtimeReq.Note,
	}
	overtime.UserID = currentUserIDUint

	updatedOvertime, err := h.service.UpdateOvertime(&overtime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update overtime"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedOvertime})
}

func (h *OvertimeHandler) DeleteOvertime(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteOvertime(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete overtime"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
