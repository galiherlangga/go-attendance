package repositories

import (
	"context"
	"time"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"gorm.io/gorm"
)

type OvertimeRepository interface {
	GetOvertimeList(userID uint, pagination utils.Pagination) ([]*models.Overtime, int64, error)
	GetOvertimeByID(id uint) (*models.Overtime, error)
	GetOvertimeByUserAndDate(userID uint, date time.Time) (*models.Overtime, error)
	CreateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error)
	UpdateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error)
	DeleteOvertime(id uint) error
	CountOvertimeHours(userID uint, startDate, endDate string) (float64, error)
}

type overtimeRepository struct {
	db *gorm.DB
}

func NewOvertimeRepository(db *gorm.DB) OvertimeRepository {
	return &overtimeRepository{
		db: db,
	}
}

func (r *overtimeRepository) GetOvertimeList(userID uint, pagination utils.Pagination) ([]*models.Overtime, int64, error) {
	var overtimeList []*models.Overtime
	var total int64

	query := r.db.Model(&models.Overtime{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((pagination.Page - 1) * pagination.Limit).
		Limit(pagination.Limit).
		Find(&overtimeList).Error; err != nil {
		return nil, 0, err
	}

	return overtimeList, total, nil
}

func (r *overtimeRepository) GetOvertimeByID(id uint) (*models.Overtime, error) {
	var overtime models.Overtime
	if err := r.db.First(&overtime, id).Error; err != nil {
		return nil, err
	}
	return &overtime, nil
}

func (r *overtimeRepository) GetOvertimeByUserAndDate(userID uint, date time.Time) (*models.Overtime, error) {
	var overtime models.Overtime
	// convert date to something like 2025-06-11
	formatedDate := date.Format("2006-01-02")
	if err := r.db.Where("user_id = ? AND date = ?", userID, formatedDate).First(&overtime).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No record found
		}
		return nil, err // Other errors
	}
	return &overtime, nil
}

func (r *overtimeRepository) CreateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error) {
	if err := r.db.WithContext(ctx).Create(overtime).Error; err != nil {
		return nil, err
	}
	return overtime, nil
}

func (r *overtimeRepository) UpdateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error) {
	if err := r.db.WithContext(ctx).Save(overtime).Error; err != nil {
		return nil, err
	}
	return overtime, nil
}

func (r *overtimeRepository) DeleteOvertime(id uint) error {
	if err := r.db.Delete(&models.Overtime{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *overtimeRepository) CountOvertimeHours(userID uint, startDate, endDate string) (float64, error) {
	var totalHours float64

	query := r.db.Model(&models.Overtime{}).
		Select("SUM(hours)").
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate)

	if err := query.Scan(&totalHours).Error; err != nil {
		return 0, err
	}

	return totalHours, nil
}