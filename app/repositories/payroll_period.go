package repositories

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"gorm.io/gorm"
)

type PayrollPeriodRepository interface {
	FindAll(pagination utils.Pagination) ([]*models.PayrollPeriod, int64, error)
	FindByID(id uint) (*models.PayrollPeriod, error)
	Create(period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	Update(period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	Delete(id uint) error
}

type payrollPeriodRepository struct {
	db *gorm.DB
}

func NewPayrollPeriodRepository(db *gorm.DB) PayrollPeriodRepository {
	return &payrollPeriodRepository{
		db: db,
	}
}

func (r *payrollPeriodRepository) FindAll(pagination utils.Pagination) ([]*models.PayrollPeriod, int64, error) {
	var periods []*models.PayrollPeriod
	var total int64
	
	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Model(&models.PayrollPeriod{})
	
	query.Count(&total)
	
	err := query.
		Limit(pagination.Limit).
		Offset(offset).
		Order("start_date DESC").
		Find(&periods).Error
	if err != nil {
		return nil, 0, err
	}
	return periods, total, nil
}

func (r *payrollPeriodRepository) FindByID(id uint) (*models.PayrollPeriod, error) {
	period := &models.PayrollPeriod{}
	if err := r.db.First(period, id).Error; err != nil {
		return nil, err // Other error
	}
	return period, nil
}

func (r *payrollPeriodRepository) Create(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	if err := r.db.Create(period).Error; err != nil {
		return nil, err // Other error
	}
	return period, nil
}

func (r *payrollPeriodRepository) Update(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	if err := r.db.Save(period).Error; err != nil {
		return nil, err // Other error
	}
	return period, nil
}

func (r *payrollPeriodRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.PayrollPeriod{}, id).Error; err != nil {
		return err // Other error
	}
	return nil
}