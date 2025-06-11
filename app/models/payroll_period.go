package models

import "gorm.io/gorm"

type PayrollPeriod struct {
	gorm.Model
	StartDate   string  `json:"start_date" gorm:"not null"`
	EndDate     string  `json:"end_date" gorm:"not null"`
	IsProcessed bool    `json:"is_processed" gorm:"default:false"`
	ProcessedAt *string `json:"processed_at" gorm:"default:null"`
}

type PayrollPeriodCache struct {
	PayrollPeriod []*PayrollPeriod `json:"payroll_periods"`
	Total         int64           `json:"total"`
}
