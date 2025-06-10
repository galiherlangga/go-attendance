package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null;size:100" json:"name"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password string `gorm:"not null;size:100" json:"password"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID;references:ID" json:"role" readonly:"true"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}
