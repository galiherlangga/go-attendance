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
	Date   time.Time `json:"date" binding:"required" format:"2006-01-02" example:"2025-06-11T00:00:00Z"`
	Amount float64   `json:"amount" binding:"required,min=0" example:"250000"`
	Note   *string   `json:"note" binding:"omitempty,max=255" example:"Client lunch and parking"`
}

type ReimbursementResponse struct {
	ID        uint       `json:"id" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2023-01-02T12:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"2023-01-10T00:00:00Z"`
	UserID    uint       `json:"user_id" example:"101"`
	Date      time.Time  `json:"date" example:"2023-06-01"`
	Amount    float64    `json:"amount" example:"250000"`
	Note      *string    `json:"note" example:"Client lunch and parking"`
}
