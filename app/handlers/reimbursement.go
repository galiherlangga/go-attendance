package handlers

import (
	"net/http"
	"strconv"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ReimbursementHandler struct {
	service services.ReimbursementService
	userService services.UserService
}

func NewReimbursementHandler(service services.ReimbursementService, userService services.UserService) *ReimbursementHandler {
	return &ReimbursementHandler{
		service: service,
		userService: userService,
	}
}

// GetReimbursementList godoc
// @Summary      Get reimbursement list
// @Description  Retrieves a list of reimbursement records for a specific user or all users if admin. Supports pagination.
// @Tags         reimbursement
// @Accept       json
// @Produce      json
// @Param        user_id   query     int  false  "User ID to filter reimbursement records"  default(0)
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        limit     query     int  false  "Number of items per page"  default(10)
// @Success      200       {object}  map[string]interface{}  "List of reimbursement records"
// @Failure      400       {object}  map[string]string        "Invalid input"
// @Failure      403       {object}  map[string]string        "Forbidden access"
// @Failure      500       {object}  map[string]string        "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /reimbursements [get]
func (h *ReimbursementHandler) GetReimbursementList(ctx *gin.Context) {
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
	// Check if the current user is not an admin and is trying to access another user's reimbursement
	if !isAdmin && uint(userID) != currentUserIDUint {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to access this user's reimbursement"})
		return
	}

	pagination := utils.GetPagination(ctx)

	reimbursements, total, err := h.service.GetReimbursementList(uint(userID), pagination)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "failed to get reimbursement list"})
		return
	}

	ctx.JSON(200, gin.H{
		"data":      reimbursements,
		"total":     total,
		"page":      pagination.Page,
		"limit":     pagination.Limit,
		"last_page": (total + int64(pagination.Limit) - 1) / int64(pagination.Limit),
	})
}


// GetReimbursementByID godoc
// @Summary      Get reimbursement by ID
// @Description  Retrieves a specific reimbursement record by its ID. Admins can access any record, while users can only access their own.
// @Tags         reimbursement
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Reimbursement ID"
// @Success      200    {object}  models.ReimbursementResponse  "Reimbursement record"
// @Failure      400    {object}  map[string]string    "Invalid ID"
// @Failure      404    {object}  map[string]string    "Reimbursement not found"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /reimbursements/{id} [get]
func (h *ReimbursementHandler) GetReimbursementByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement ID"})
		return
	}

	reimbursement, err := h.service.GetReimbursementByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reimbursement not found"})
		return
	}

	ctx.JSON(http.StatusOK, reimbursement)
}


// CreateReimbursement godoc
// @Summary      Create reimbursement
// @Description  Creates a new reimbursement record for the current user. The user must be authenticated.
// @Tags         reimbursement
// @Accept       json
// @Produce      json
// @Param        body   body      models.ReimbursementRequest  true  "Reimbursement payload"
// @Success      201    {object}  models.ReimbursementResponse  "Created reimbursement record"
// @Failure      400    {object}  map[string]string  "Invalid input"
// @Failure      401    {object}  map[string]string  "Unauthorized"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /reimbursements [post]
func (h *ReimbursementHandler) CreateReimbursement(ctx *gin.Context) {
	var reimbursementReq models.ReimbursementRequest
	if err := ctx.ShouldBindJSON(&reimbursementReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

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

	reimbursement := models.Reimbursement{
		UserID: currentUserIDUint,
		Date: reimbursementReq.Date,
		Amount: reimbursementReq.Amount,
		Note:   reimbursementReq.Note,
	}
	newReimbursement, err := h.service.SubmitReimbursement(&reimbursement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reimbursement"})
		return
	}

	ctx.JSON(http.StatusCreated, newReimbursement)
}


// UpdateReimbursement godoc
// @Summary      Update reimbursement
// @Description  Updates an existing reimbursement record by its ID. Only the user who created the reimbursement or an admin can update it.
// @Tags         reimbursement
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Reimbursement ID"
// @Param        body   body      models.ReimbursementRequest  true  "Reimbursement payload"
// @Success      200    {object}  models.ReimbursementResponse  "Updated reimbursement record"
// @Failure      400    {object}  map[string]string  "Invalid input"
// @Failure      404    {object}  map[string]string  "Reimbursement not found"
// @Failure      403    {object}  map[string]string  "Forbidden access"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /reimbursements/{id} [put]
func (h *ReimbursementHandler) UpdateReimbursement(ctx *gin.Context) {
	var reimbursementReq models.ReimbursementRequest
	if err := ctx.ShouldBindJSON(&reimbursementReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement ID"})
		return
	}

	reimbursement, err := h.service.GetReimbursementByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "reimbursement not found"})
		return
	}
	reimbursement.Date = reimbursementReq.Date
	reimbursement.Amount = reimbursementReq.Amount
	reimbursement.Note = reimbursementReq.Note
	updatedReimbursement, err := h.service.UpdateReimbursement(reimbursement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update reimbursement"})
		return
	}

	ctx.JSON(http.StatusOK, updatedReimbursement)
}


// DeleteReimbursement godoc
// @Summary      Delete reimbursement
// @Description  Deletes an existing reimbursement record by its ID. Only the user who created the reimbursement or an admin can delete it.
// @Tags         reimbursement
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Reimbursement ID"
// @Success      204    "No Content"
// @Failure      400    {object}  map[string]string  "Invalid ID"
// @Failure      403    {object}  map[string]string  "Forbidden access"
// @Failure      404    {object}  map[string]string  "Reimbursement not found"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /reimbursements/{id} [delete]
func (h *ReimbursementHandler) DeleteReimbursement(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reimbursement ID"})
		return
	}

	err = h.service.DeleteReimbursement(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete reimbursement"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}