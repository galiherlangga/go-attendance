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

// GetOvertimeList godoc
// @Summary      Get overtime list
// @Description  Retrieves a list of overtime records for a specific user or all users if admin. Supports pagination.
// @Tags         overtime
// @Accept       json
// @Produce      json
// @Param        user_id   query     int  false  "User ID to filter overtime records"  default(0)
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        limit     query     int  false  "Number of items per page"  default(10)
// @Success      200       {object}  map[string]interface{}  "List of overtime records"
// @Failure      400       {object}  map[string]string        "Invalid input"
// @Failure      403       {object}  map[string]string        "Forbidden access"
// @Failure      500       {object}  map[string]string        "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /overtimes [get]
func (h *OvertimeHandler) GetOvertimeList(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user access"})
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
	// Check if the current user is not an admin and is trying to access another user's overtime
	if !isAdmin && uint(userID) != currentUserIDUint {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to access this user's overtime"})
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


// GetOvertimeByID godoc
// @Summary      Get overtime by ID
// @Description  Retrieves a specific overtime record by its ID. Admins can access any record, while users can only access their own.
// @Tags         overtime
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Overtime ID"
// @Success      200    {object}  map[string]interface{}  "Overtime record"
// @Failure      400    {object}  map[string]string        "Invalid ID"
// @Failure      404    {object}  map[string]string        "Overtime not found"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /overtimes/{id} [get]
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


// CreateOvertime godoc
// @Summary      Create overtime
// @Description  Creates a new overtime record for the authenticated user.
// @Tags         overtime
// @Accept       json
// @Produce      json
// @Param        body  body      models.OvertimeRequest  true  "Overtime request payload"
// @Success      201   {object}  models.OvertimeResponse  "Created overtime record"
// @Failure      400   {object}  map[string]string  "Invalid input"
// @Failure      401   {object}  map[string]string  "Unauthorized"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /overtimes [post]
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


// UpdateOvertime godoc
// @Summary      Update overtime
// @Description  Updates an existing overtime record for the authenticated user.
// @Tags         overtime
// @Accept       json
// @Produce      json
// @Param        body  body      models.OvertimeRequest  true  "Overtime request payload"
// @Success      200   {object}  models.OvertimeResponse  "Updated overtime record"
// @Failure      400   {object}  map[string]string  "Invalid input"
// @Failure      401   {object}  map[string]string  "Unauthorized"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /overtimes [put]
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


// DeleteOvertime godoc
// @Summary      Delete overtime
// @Description  Deletes an existing overtime record by its ID. Admins can delete any record, while users can only delete their own.
// @Tags         overtime
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Overtime ID"
// @Success      204    "No Content"
// @Failure      400    {object}  map[string]string  "Invalid ID"
// @Failure      403    {object}  map[string]string  "Forbidden access"
// @Failure      404    {object}  map[string]string  "Overtime not found"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /overtimes/{id} [delete]
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
