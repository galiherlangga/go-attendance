package repositories

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Preload("Role").Where("email = ?", email).First(user).Error; err != nil {
		return nil, err // Other error
	}
	return user, nil
}