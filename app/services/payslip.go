package services

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
)

type PayslipService interface {
	GeneratePayslip(ctx context.Context, userID uint, periodID uint, monthlySalary float64) error
	GetSummary(periodID uint) (*models.PayslipSummary, float64, error)
	GetPayslipByUserAndPeriod(userID uint, periodID uint) (*models.Payslip, error)
}

type payslipService struct {
	repo repositories.PayslipRepository
	attendanceRepo repositories.AttendanceRepository
	overtimeRepo repositories.OvertimeRepository
	reimbursementRepo repositories.ReimbursementRepository
	periodRepo repositories.PayrollPeriodRepository
	userRepo repositories.UserRepository
}

func NewPayslipService(
	repo repositories.PayslipRepository,
	attendanceRepo repositories.AttendanceRepository,
	overtimeRepo repositories.OvertimeRepository,
	reimbursementRepo repositories.ReimbursementRepository,
	periodRepo repositories.PayrollPeriodRepository,
	userRepo repositories.UserRepository) PayslipService {
	return &payslipService{
		repo: repo,
		attendanceRepo: attendanceRepo,
		overtimeRepo: overtimeRepo,
		reimbursementRepo: reimbursementRepo,
		periodRepo: periodRepo,
	}
}

func (s *payslipService) GeneratePayslip(ctx context.Context, userID uint, periodID uint, monthlySalary float64) error {
	existing, _ := s.repo.GetByUserAndPeriod(userID, periodID)
	if existing != nil {
		return errors.New("payslip already generated for this period")
	}
	period, err := s.periodRepo.FindByID(periodID)
	if err != nil {
		return err
	}
	
	start := period.StartDate.Format("2006-01-02")
	end := period.EndDate.Format("2006-01-02")
	workdays := utils.CountWorkingDays(period.StartDate, period.EndDate)
	attended, _ := s.attendanceRepo.CountWorkingDays(userID, start, end)
	overtimeHours, _ := s.overtimeRepo.CountOvertimeHours(userID, start, end)
	reimbursements, _ := s.reimbursementRepo.SumReimbursement(userID, start, end)
	
	dailySalary := monthlySalary / float64(workdays)
	dailySalary = math.Round(dailySalary * 100) / 100
	overtimePay := overtimeHours * (dailySalary / 8) * 2
	overtimePay = math.Round(overtimePay * 100) / 100
	total := float64(attended) * dailySalary + overtimePay + reimbursements
	total = math.Round(total * 100) / 100
	
	now := time.Now()
	payslip := &models.Payslip{
		UserID: userID,
		PayrollPeriodID: periodID,
		GeneratedAt: now,
		AttendanceDays: int(attended),
		AttendanceEarnings: dailySalary,
		OvertimeHours: overtimeHours,
		OvertimeEarnings: overtimePay,
		TotalReimbursement: reimbursements,
		TakeHomePay: total,
	}
	if err := s.repo.Create(ctx, payslip); err != nil {
		return err
	}
	return nil
}

func (s *payslipService) GetSummary(periodID uint) (*models.PayslipSummary, float64, error) {
	payslips, err := s.repo.GetByPeriod(periodID)
	if err != nil {
		return nil, 0, err
	}

	var summaryItems []models.PayslipSummaryItem
	total := 0.0

	for _, payslip := range payslips {
		summaryItems = append(summaryItems, models.PayslipSummaryItem{
			UserID:      payslip.UserID,
			Name:        payslip.User.Name,
			TakeHomePay: payslip.TakeHomePay,
		})
		total += payslip.TakeHomePay
	}

	total = math.Round(total * 100) / 100

	return &models.PayslipSummary{Items: summaryItems}, total, nil
}

func (s *payslipService) GetPayslipByUserAndPeriod(userID uint, periodID uint) (*models.Payslip, error) {
	payslip, err := s.repo.GetByUserAndPeriod(userID, periodID)
	if err != nil {
		return nil, err
	}
	if payslip == nil {
		return nil, errors.New("payslip not found")
	}
	return payslip, nil
}