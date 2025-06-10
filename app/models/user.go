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