package services

import (
	"context"
	"errors"
	"time"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type PayrollPeriodService interface {
	GetPayrollPeriodList(pagination utils.Pagination) ([]*models.PayrollPeriod, int64, error)
	GetPayrollPeriodByID(id uint) (*models.PayrollPeriod, error)
	CreatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	UpdatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	DeletePayrollPeriod(id uint) error
}

type payrollPeriodService struct {
	repo  repositories.PayrollPeriodRepository
	cache *redis.Client
}

func NewPayrollPeriodService(repo repositories.PayrollPeriodRepository, cache *redis.Client) PayrollPeriodService {
	return &payrollPeriodService{
		repo:  repo,
		cache: cache,
	}
}

func (s *payrollPeriodService) GetPayrollPeriodList(pagination utils.Pagination) ([]*models.PayrollPeriod, int64, error) {
	ctx := context.Background()
	cacheKey := utils.BuildKey("payroll", pagination.Page, pagination.Limit)

	if cached, err := utils.GetCache[models.PayrollPeriodCache](ctx, s.cache, cacheKey); err == nil {
		return cached.PayrollPeriod, cached.Total, nil
	}

	// Fallback to DB
	payrollPeriods, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, &models.PayrollPeriodCache{
		PayrollPeriod: payrollPeriods,
		Total:         total,
	}, 10*time.Minute)
	return payrollPeriods, total, nil
}

func (s *payrollPeriodService) GetPayrollPeriodByID(id uint) (*models.PayrollPeriod, error) {
	ctx := context.Background()
	cacheKey := utils.BuildKey("payroll", id)

	if cached, err := utils.GetCache[models.PayrollPeriod](ctx, s.cache, cacheKey); err == nil {
		return cached, nil
	}

	// Fallback to DB
	period, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, period, 10*time.Minute)
	return period, nil
}

func (s *payrollPeriodService) CreatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	if period == nil {
		return nil, errors.New("payroll period cannot be nil")
	}

	createdPeriod, err := s.repo.Create(period)
	if err != nil {
		return nil, err
	}

	return createdPeriod, nil
}

func (s *payrollPeriodService) UpdatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	if period == nil {
		return nil, errors.New("payroll period cannot be nil")
	}

	updatedPeriod, err := s.repo.Update(period)
	if err != nil {
		return nil, err
	}
	
	ctx := context.Background()
	cacheKey := utils.BuildKey("payroll", period.ID)
	s.cache.Del(ctx, cacheKey) // Remove cache if exists

	return updatedPeriod, nil
}

func (s *payrollPeriodService) DeletePayrollPeriod(id uint) error {
	if id == 0 {
		return errors.New("invalid payroll period ID")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	ctx := context.Background()
	cacheKey := utils.BuildKey("payroll", id)
	s.cache.Del(ctx, cacheKey) // Invalidate cache after deletion

	return nil
}
