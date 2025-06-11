package models

import (
	"time"

	"gorm.io/gorm"
)

type Reimbursement struct {
	gorm.Model
	UserID uint      `json:"user_id" gorm:"not null,uniqueIndex:idx_user_date"`
	Date   time.Time `json:"date" gorm:"type:DATE;not null,uniqueIndex:idx_user_date"`
	Amount float64   `json:"amount" gorm:"not null"`
	Note   *string   `json:"note" gorm:"type:text"`
}

type ReimbursementCache struct {
	ReimbursementList []*Reimbursement `json:"reimbursements"`
	Total             int64            `json:"total"`
}

type ReimbursementRequest struct {
	Date   time.Time `json:"date" binding:"required" format:"2006-01-02"`
	Amount float64   `json:"amount" binding:"required,min=0"`
	Note   *string   `json:"note" binding:"omitempty,max=255"`
}
