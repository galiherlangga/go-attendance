package repositories

import (
	"context"

	"github.com/galiherlangga/go-attendance/app/models"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	GetAttendanceList(userID uint, startDate, endDate string) ([]*models.Attendance, error)
	GetAttendanceByID(id uint) (*models.Attendance, error)
	GetAttendanceByUserAndDate(userID uint, date string) (*models.Attendance, error)
	CreateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	UpdateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	DeleteAttendance(id uint) error
	CountWorkingDays(userID uint, startDate, endDate string) (int64, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{
		db: db,
	}
}

func (r *attendanceRepository) GetAttendanceList(userID uint, startDate, endDate string) ([]*models.Attendance, error) {
	var attendances []*models.Attendance
	query := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate)
	if err := query.Find(&attendances).Error; err != nil {
		return nil, err
	}
	return attendances, nil
}

func (r *attendanceRepository) GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := r.db.First(&attendance, id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) GetAttendanceByUserAndDate(userID uint, date string) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := r.db.Where("user_id = ? AND date = ?", userID, date).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) CreateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error) {
	if err := r.db.WithContext(ctx).Create(attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *attendanceRepository) UpdateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error) {
	if err := r.db.WithContext(ctx).Save(attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *attendanceRepository) DeleteAttendance(id uint) error {
	if err := r.db.Delete(&models.Attendance{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *attendanceRepository) CountWorkingDays(userID uint, startDate, endDate string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).
		Where("check_in IS NOT NULL AND check_out IS NOT NULL")
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}