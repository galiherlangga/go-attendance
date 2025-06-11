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