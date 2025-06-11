package repositories

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"gorm.io/gorm"
)

type ReimbursementRepository interface {
	GetReimbursementList(userID uint, pagination utils.Pagination) ([]*models.Reimbursement, int64, error)
	GetReimbursementByID(id uint) (*models.Reimbursement, error)
	CreateReimbursement(reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	UpdateReimbursement(reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	DeleteReimbursement(id uint) error
}

type reimbursementRepository struct {
	db *gorm.DB
}

func NewReimbursementRepository(db *gorm.DB) ReimbursementRepository {
	return &reimbursementRepository{
		db: db,
	}
}

func (r *reimbursementRepository) GetReimbursementList(userID uint, pagination utils.Pagination) ([]*models.Reimbursement, int64, error) {
	var reimbursements []*models.Reimbursement
	var total int64

	query := r.db.Model(&models.Reimbursement{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&reimbursements).Error; err != nil {
		return nil, 0, err
	}
	return reimbursements, total, nil
}

func (r *reimbursementRepository) GetReimbursementByID(id uint) (*models.Reimbursement, error) {
	var reimbursement models.Reimbursement
	if err := r.db.First(&reimbursement, id).Error; err != nil {
		return nil, err
	}
	return &reimbursement, nil
}

func (r *reimbursementRepository) CreateReimbursement(reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	if err := r.db.Create(reimbursement).Error; err != nil {
		return nil, err
	}
	return reimbursement, nil
}

func (r *reimbursementRepository) UpdateReimbursement(reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	if err := r.db.Save(reimbursement).Error; err != nil {
		return nil, err
	}
	return reimbursement, nil
}

func (r *reimbursementRepository) DeleteReimbursement(id uint) error {
	if err := r.db.Delete(&models.Reimbursement{}, id).Error; err != nil {
		return err
	}
	return nil
}
