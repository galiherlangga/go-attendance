package handlers

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/services"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PayrollPeriodHandler struct {
	service services.PayrollPeriodService
}

func NewPayrollPeriodHandler(service services.PayrollPeriodService) *PayrollPeriodHandler {
	return &PayrollPeriodHandler{
		service: service,
	}
}

// GetPayrollPeriodList godoc
// @Summary      Get list of payroll periods
// @Description  Retrieves a paginated list of payroll periods. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Number of items per page"  default(10)
// @Success      200    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods [get]
func (h *PayrollPeriodHandler) GetPayrollPeriodList(ctx *gin.Context) {
	pagination := utils.GetPagination(ctx)

	payrollPeriods, total, err := h.service.GetPayrollPeriodList(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payroll periods"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":      payrollPeriods,
		"total":     total,
		"page":      pagination.Page,
		"limit":     pagination.Limit,
		"totalPage": int(math.Ceil(float64(total) / float64(pagination.Limit))),
	})
}

// GetPayrollPeriodByID godoc
// @Summary      Get payroll period by ID
// @Description  Retrieves a specific payroll period by its ID. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Payroll Period ID"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods/{id} [get]
func (h *PayrollPeriodHandler) GetPayrollPeriodByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	payrollPeriod, err := h.service.GetPayrollPeriodByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Payroll period not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": payrollPeriod})
}

// CreatePayrollPeriod godoc
// @Summary      Create a payroll period
// @Description  Creates a new payroll period. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        body   body      models.PayrollPeriodExample  true  "Payroll Period payload"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods [post]
func (h *PayrollPeriodHandler) CreatePayrollPeriod(ctx *gin.Context) {
	var period models.PayrollPeriod
	if err := ctx.ShouldBindJSON(&period); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	requestID := ctx.GetString("request_id")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user_id", userID))
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_id", requestID))

	createdPeriod, err := h.service.CreatePayrollPeriod(ctx, &period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payroll period"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdPeriod})
}

// UpdatePayrollPeriod godoc
// @Summary      Update a payroll period
// @Description  Updates an existing payroll period. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Payroll Period ID"
// @Param        body   body      models.PayrollPeriod  true  "Payroll Period payload"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods/{id} [put]
func (h *PayrollPeriodHandler) UpdatePayrollPeriod(ctx *gin.Context) {
	var period models.PayrollPeriod
	if err := ctx.ShouldBindJSON(&period); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	periodID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	period.ID = uint(periodID)

	userID := ctx.GetUint("user_id")
	requestID := ctx.GetString("request_id")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user_id", userID))
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "request_id", requestID))

	updatedPeriod, err := h.service.UpdatePayrollPeriod(ctx, &period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payroll period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPeriod})
}

// DeletePayrollPeriod godoc
// @Summary      Delete a payroll period
// @Description  Deletes a payroll period by its ID. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Payroll Period ID"
// @Success      204    "No Content"
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods/{id} [delete]
func (h *PayrollPeriodHandler) DeletePayrollPeriod(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeletePayrollPeriod(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payroll period"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// RunPayrollPeriod godoc
// @Summary      Run payroll period
// @Description  Runs the payroll calculations for a specific payroll period. Admin only.
// @Tags         payroll
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Payroll Period ID"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     CookieAuth
// @Security     BearerAuth
// @Router       /payroll-periods/{id}/run-payroll [post]
func (h *PayrollPeriodHandler) RunPayrollPeriod(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := ctx.GetUint("user_id")
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "user_id", userID))

	if err := h.service.RunPayroll(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run payroll period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payroll period successfully run"})
}
