package models

import (
	"time"

	"gorm.io/gorm"
)

type PayrollPeriod struct {
	gorm.Model
	StartDate   time.Time `json:"start_date" gorm:"not null"`
	EndDate     time.Time `json:"end_date" gorm:"not null"`
	IsProcessed bool      `json:"is_processed" gorm:"default:false"`
	ProcessedAt *string   `json:"processed_at" gorm:"default:null"`
}

type PayrollPeriodCache struct {
	PayrollPeriod []*PayrollPeriod `json:"payroll_periods"`
	Total         int64            `json:"total"`
}

// PayrollPeriodExample is used only for Swagger documentation.
// @Description Sample payload for creating a payroll period
type PayrollPeriodExample struct {
	StartDate string `json:"start_date" example:"2023-01-01T00:00:00Z"`
	EndDate   string `json:"end_date" example:"2023-01-15T00:00:00Z"`
}