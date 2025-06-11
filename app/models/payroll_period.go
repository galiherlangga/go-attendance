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
