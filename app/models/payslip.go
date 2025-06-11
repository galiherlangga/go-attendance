package models

import (
	"time"

	"gorm.io/gorm"
)

type Payslip struct {
	gorm.Model
	UserID             uint      `json:"user_id" gorm:"not null,uniqueIndex:idx_user_period"`
	PayrollPeriodID    uint      `json:"payroll_period_id" gorm:"not null,uniqueIndex:idx_user_period"`
	GeneratedAt        time.Time `json:"generated_at" gorm:"default:null"`
	AttendanceDays     int       `json:"attendance_days" gorm:"not null"`
	AttendanceEarnings float64   `json:"attendance_earnings" gorm:"not null"`
	OvertimeHours      float64   `json:"overtime_hours" gorm:"not null"`
	OvertimeEarnings   float64   `json:"overtime_earnings" gorm:"not null"`
	TotalReimbursement float64   `json:"total_reimbursement" gorm:"not null"`
	TakeHomePay        float64   `json:"take_home_pay" gorm:"not null"`
	User               User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" readonly:"true"`
}

type PayslipSummary struct {
	Items []PayslipSummaryItem `json:"items"`
}

type PayslipSummaryItem struct {
	UserID      uint    `json:"user_id"`
	Name        string  `json:"name"`
	TakeHomePay float64 `json:"take_home_pay"`
}
