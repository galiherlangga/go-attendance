package models

import (
	"time"

	"gorm.io/gorm"
)

type Overtime struct {
	gorm.Model
	UserID uint      `json:"user_id" gorm:"not null,uniqueIndex:idx_user_date"`
	Date   time.Time `json:"date" gorm:"type:DATE;not null,uniqueIndex:idx_user_date"`
	Hours  int       `json:"hours" gorm:"not null"`
	Note   *string   `json:"note" gorm:"size:255"`
}

type OvertimeCache struct {
	OvertimeList []*Overtime `json:"overtime"`
	Total        int64       `json:"total"`
}

type OvertimeRequest struct {
	Date  time.Time `json:"date" binding:"required" format:"2006-01-02"`
	Hours int       `json:"hours" binding:"required,min=1"`
	Note  *string   `json:"note" binding:"omitempty,max=255"`
}

// OvertimeResponse is used for Swagger documentation
type OvertimeResponse struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uint       `json:"user_id"`
	Date      time.Time  `json:"date"`
	Hours     int        `json:"hours"`
	Note      *string    `json:"note"`
}
