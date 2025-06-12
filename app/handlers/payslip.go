package handlers

import (
	"net/http"
	"strconv"

	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/gin-gonic/gin"
)

type PayslipHandler struct {
	service services.PayslipService
	userService services.UserService
}

func NewPayslipHandler(service services.PayslipService, userService services.UserService) *PayslipHandler {
	return &PayslipHandler{
		service: service,
		userService: userService,
	}
}

// GetPayslipList godoc
// @Summary      Get payslip list
// @Description  Retrieves a list of payslips for a specific user or all users if admin. Supports pagination.
// @Tags         payslip
// @Accept       json
// @Produce      json
// @Param        user_id   query     int  false  "User ID to filter payslip records"  default(0)
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        limit     query     int  false  "Number of items per page"  default(10)
// @Success      200       {object}  map[string]interface{}  "List of payslip records"
// @Failure      400       {object}  map[string]string        "Invalid input"
// @Failure      403       {object}  map[string]string        "Forbidden access"
// @Failure      500       {object}  map[string]string        "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payslips [get]
func (h *PayslipHandler) GetPayslipByUserAndPeriod(ctx *gin.Context) {
	currentUserID, exists := ctx.Get("user_id")
	userIDParam := ctx.DefaultQuery("user_id", "0")
	periodIDParam := ctx.Param("period_id")
	// Check if the user is authenticated
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

	periodID, err := strconv.ParseUint(periodIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}

	payslip, err := h.service.GetPayslipByUserAndPeriod(uint(userID), uint(periodID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payslip: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payslip)
}

// GetPayslipSummary godoc
// @Summary      Get payslip summary
// @Description  Retrieves a summary of payslips for a specific payroll period.
// @Tags         payslip
// @Accept       json
// @Produce      json
// @Param        period_id  path      int  true  "Payroll Period ID"
// @Success      200         {object}  map[string]interface{}  "Payslip summary"
// @Failure      400         {object}  map[string]string        "Invalid input"
// @Failure      500         {object}  map[string]string        "Internal server error"
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payslips/summary/{period_id} [get]
func (h *PayslipHandler) GetPayslipSummary(ctx *gin.Context) {
	periodIDParam := ctx.Param("period_id")
	periodID, err := strconv.ParseUint(periodIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_id"})
		return
	}

	summary, total, err := h.service.GetSummary(uint(periodID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payslip summary: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"total":   total,
	})
}
