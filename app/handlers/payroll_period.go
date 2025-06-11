package handlers

import (
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

func (h *PayrollPeriodHandler) CreatePayrollPeriod(ctx *gin.Context) {
	var period models.PayrollPeriod
	if err := ctx.ShouldBindJSON(&period); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	createdPeriod, err := h.service.CreatePayrollPeriod(&period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payroll period"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": createdPeriod})
}

func (h *PayrollPeriodHandler) UpdatePayrollPeriod(ctx *gin.Context) {
	var period models.PayrollPeriod
	if err := ctx.ShouldBindJSON(&period); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedPeriod, err := h.service.UpdatePayrollPeriod(&period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payroll period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPeriod})
}

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

func (h *PayrollPeriodHandler) RunPayrollPeriod(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.RunPayroll(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run payroll period"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Payroll period successfully run"})
}