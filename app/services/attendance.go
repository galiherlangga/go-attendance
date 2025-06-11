package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"gorm.io/gorm"
)

type AttendanceService interface {
	GetAttendanceList(userID uint, startDate, endDate string) ([]*models.Attendance, error)
	GetAttendanceByID(id uint) (*models.Attendance, error)
	CheckIn(userID uint) (*models.Attendance, error)
	CheckOut(userID uint) (*models.Attendance, error)
}

type attendanceService struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceService(repo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{
		repo: repo,
	}
}

func (s *attendanceService) GetAttendanceList(userID uint, startDate, endDate string) ([]*models.Attendance, error) {
	return s.repo.GetAttendanceList(userID, startDate, endDate)
}

func (s *attendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	return s.repo.GetAttendanceByID(id)
}

func (s *attendanceService) CheckIn(userID uint) (*models.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	if utils.IsWeekend(today) {
		return nil, errors.New("cannot submit attendance on weekends")
	}

	attDate := today.Format("2006-01-02")
	att, err := s.repo.GetAttendanceByUserAndDate(userID, attDate)
	fmt.Println(att)
	now := time.Now()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		att = &models.Attendance{
			UserID:  userID,
			Date:    today,
			CheckIn: &now,
		}
		return s.repo.CreateAttendance(att)
	}

	return nil, errors.New("attendance already exists for today")
}

func (s *attendanceService) CheckOut(userID uint) (*models.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	if utils.IsWeekend(today) {
		return nil, errors.New("cannot submit attendance on weekends")
	}

	attDate := today.Format("2006-01-02")
	att, err := s.repo.GetAttendanceByUserAndDate(userID, attDate)
	if err != nil {
		return nil, err
	}

	if att.CheckIn == nil {
		return nil, errors.New("check-in must be done before check-out")
	}
	
	if att.CheckOut != nil {
		return nil, errors.New("check-out already done for today")
	}

	now := time.Now()
	att.CheckOut = &now
	return s.repo.UpdateAttendance(att)
}
