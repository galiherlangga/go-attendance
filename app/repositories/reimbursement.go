package repositories

import (
	"context"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"gorm.io/gorm"
)

type ReimbursementRepository interface {
	GetReimbursementList(userID uint, pagination utils.Pagination) ([]*models.Reimbursement, int64, error)
	GetReimbursementByID(id uint) (*models.Reimbursement, error)
	CreateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	UpdateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	DeleteReimbursement(id uint) error
	SumReimbursement(userID uint, startDate, endDate string) (float64, error)
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

func (r *reimbursementRepository) CreateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	if err := r.db.WithContext(ctx).Create(reimbursement).Error; err != nil {
		return nil, err
	}
	return reimbursement, nil
}

func (r *reimbursementRepository) UpdateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	if err := r.db.WithContext(ctx).Save(reimbursement).Error; err != nil {
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

func (r *reimbursementRepository) SumReimbursement(userID uint, startDate, endDate string) (float64, error) {
	var total float64
	if err := r.db.Model(&models.Reimbursement{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Select("SUM(amount)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}