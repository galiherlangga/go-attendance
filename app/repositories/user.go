package repositories

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	GetAllEmployee(offset int, limit int) ([]*models.User, error)
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

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Preload("Role").Where("id = ?", id).First(user).Error; err != nil {
		return nil, err // Other error
	}
	return user, nil
}

func (r *userRepository) GetAllEmployee(offset int, limit int) ([]*models.User, error) {
	var users []*models.User
	employeeRole := &models.Role{}
	err := r.db.Model(&models.Role{}).Where("name = 'user'").First(employeeRole).Error; if err != nil {
		return nil, err
	}

	query := r.db.Model(&models.User{}).Where("role_id = ?",employeeRole.ID)

	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
