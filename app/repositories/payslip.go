package repositories

import (
	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

type PayslipRepository interface {
	Create(payslip *models.Payslip) error
	GetByUserAndPeriod(userID uint, periodID uint) (*models.Payslip, error)
	GetByID(id uint) (*models.Payslip, error)
	GetByPeriod(periodID uint) ([]*models.Payslip, error)
	Exists(userID uint, periodID uint) (bool, error)
}

type payslipRepository struct {
	db *gorm.DB
}

func NewPayslipRepository(db *gorm.DB) PayslipRepository {
	return &payslipRepository{
		db: db,
	}
}

func (r *payslipRepository) Create(payslip *models.Payslip) error {
	return r.db.Create(payslip).Error
}

func (r *payslipRepository) GetByUserAndPeriod(userID uint, periodID uint) (*models.Payslip, error) {
	var payslip models.Payslip
	err := r.db.Where("user_id = ? AND period_id = ?", userID, periodID).First(&payslip).Error
	if err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (r *payslipRepository) GetByID(id uint) (*models.Payslip, error) {
	var payslip models.Payslip
	err := r.db.First(&payslip, id).Error
	if err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (r *payslipRepository) GetByPeriod(periodID uint) ([]*models.Payslip, error) {
	var payslips []*models.Payslip
	err := r.db.Where("period_id = ?", periodID).Find(&payslips).Error
	if err != nil {
		return nil, err
	}
	return payslips, nil
}

func (r *payslipRepository) Exists(userID uint, periodID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Payslip{}).Where("user_id = ? AND period_id = ?", userID, periodID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


