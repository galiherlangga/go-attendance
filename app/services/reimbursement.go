package services

import (
	"context"
	"errors"

	"github.com/galiherlangga/go-attendance/app/models"
	"github.com/galiherlangga/go-attendance/app/repositories"
	"github.com/galiherlangga/go-attendance/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type ReimbursementService interface {
	GetReimbursementList(userID uint, pagination utils.Pagination) ([]*models.Reimbursement, int64, error)
	GetReimbursementByID(id uint) (*models.Reimbursement, error)
	SubmitReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	UpdateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error)
	DeleteReimbursement(id uint) error
}

type reimbursementService struct {
	repo              repositories.ReimbursementRepository
	payrollPeriodRepo repositories.PayrollPeriodRepository
	cache             *redis.Client
}

func NewReimbursementService(repo repositories.ReimbursementRepository, payrollPeriodRepo repositories.PayrollPeriodRepository, cache *redis.Client) ReimbursementService {
	return &reimbursementService{
		repo:              repo,
		payrollPeriodRepo: payrollPeriodRepo,
		cache:             cache,
	}
}

func (s *reimbursementService) GetReimbursementList(userID uint, pagination utils.Pagination) ([]*models.Reimbursement, int64, error) {
	ctx := context.Background()

	cacheKey := utils.BuildKey("reimbursement", userID, pagination.Page, pagination.Limit)

	if cached, err := utils.GetCache[models.ReimbursementCache](ctx, s.cache, cacheKey); err == nil {
		return cached.ReimbursementList, cached.Total, nil
	}

	// Fallback to DB
	reimbursements, total, err := s.repo.GetReimbursementList(userID, pagination)
	if err != nil {
		return nil, 0, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, &models.ReimbursementCache{
		ReimbursementList: reimbursements,
		Total:             total,
	}, 10)

	return reimbursements, total, nil
}

func (s *reimbursementService) GetReimbursementByID(id uint) (*models.Reimbursement, error) {
	ctx := context.Background()
	cacheKey := utils.BuildKey("reimbursement", id)

	if cached, err := utils.GetCache[models.Reimbursement](ctx, s.cache, cacheKey); err == nil {
		return cached, nil
	}

	// Fallback to DB
	reimbursement, err := s.repo.GetReimbursementByID(id)
	if err != nil {
		return nil, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, reimbursement, 10)

	return reimbursement, nil
}

func (s *reimbursementService) SubmitReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	if reimbursement.Amount <= 0 {
		return nil, errors.New("reimbursement amount must be greater than zero")
	}

	isLocked, err := s.payrollPeriodRepo.IsDateLocked(reimbursement.Date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	if isLocked {
		return nil, errors.New("reimbursement cannot be submitted for a locked payroll period")
	}

	// Create reimbursement in DB
	createdReimbursement, err := s.repo.CreateReimbursement(ctx, reimbursement)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for the user
	err = utils.DeleteCacheByPattern(ctx, s.cache, "reimbursement:*")
	if err != nil {
		return nil, err
	}

	return createdReimbursement, nil
}

func (s *reimbursementService) UpdateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) (*models.Reimbursement, error) {
	isLocked, err := s.payrollPeriodRepo.IsDateLocked(reimbursement.Date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	if isLocked {
		return nil, errors.New("reimbursement cannot be submitted for a locked payroll period")
	}

	// Update reimbursement in DB
	updatedReimbursement, err := s.repo.UpdateReimbursement(ctx, reimbursement)
	if err != nil {
		return nil, err
	}

	// Invalidate cache for the user
	err = utils.DeleteCacheByPattern(ctx, s.cache, "reimbursement:*")

	return updatedReimbursement, nil
}

func (s *reimbursementService) DeleteReimbursement(id uint) error {
	ctx := context.Background()

	// Delete reimbursement in DB
	if err := s.repo.DeleteReimbursement(id); err != nil {
		return err
	}

	// Invalidate cache for the user
	err := utils.DeleteCacheByPattern(ctx, s.cache, "reimbursement:*")
	return err
}
