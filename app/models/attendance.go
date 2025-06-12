package models

import (
	"time"
)

type Attendance struct {
	BaseModel
	UserID   uint       `json:"user_id" gorm:"not null,uniqueIndex:idx_user_date"`
	Date     time.Time  `json:"date" gorm:"type:DATE;not null,uniqueIndex:idx_user_date"`
	CheckIn  *time.Time `json:"check_in" gorm:"not null"`
	CheckOut *time.Time `json:"check_out" gorm:"nullable"`
	User     User       `gorm:"foreignKey:UserID" json:"user" readonly:"true"`
}
