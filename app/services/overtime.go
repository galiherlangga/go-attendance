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

const maxHourLimit = 3

type OvertimeService interface {
	GetOvertimeList(userID uint, pagination utils.Pagination) ([]*models.Overtime, int64, error)
	GetOvertimeByID(id uint) (*models.Overtime, error)
	SubmitOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error)
	UpdateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error)
	DeleteOvertime(id uint) error
}

type overtimeService struct {
	repo  repositories.OvertimeRepository
	cache *redis.Client
}

func NewOvertimeService(repo repositories.OvertimeRepository, cache *redis.Client) OvertimeService {
	return &overtimeService{
		repo:  repo,
		cache: cache,
	}
}

func (s *overtimeService) GetOvertimeList(userID uint, pagination utils.Pagination) ([]*models.Overtime, int64, error) {
	ctx := context.Background()
	cacheKey := utils.BuildKey("overtime", pagination.Page, pagination.Limit)

	if cached, err := utils.GetCache[models.OvertimeCache](ctx, s.cache, cacheKey); err == nil {
		return cached.OvertimeList, cached.Total, nil
	}

	// Fallback to DB
	overtimeList, total, err := s.repo.GetOvertimeList(userID, pagination)
	if err != nil {
		return nil, 0, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, &models.OvertimeCache{
		OvertimeList: overtimeList,
		Total:        total,
	}, 10*time.Minute)

	return overtimeList, total, nil
}

func (s *overtimeService) GetOvertimeByID(id uint) (*models.Overtime, error) {
	ctx := context.Background()
	cacheKey := utils.BuildKey("overtime", id)

	if cached, err := utils.GetCache[models.Overtime](ctx, s.cache, cacheKey); err == nil {
		return cached, nil
	}

	// Fallback to DB
	overtime, err := s.repo.GetOvertimeByID(id)
	if err != nil {
		return nil, err
	}

	// Save to cache
	utils.SetCache(ctx, s.cache, cacheKey, overtime, 10*time.Minute)

	return overtime, nil
}

func (s *overtimeService) SubmitOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error) {
	exists, err := s.repo.GetOvertimeByUserAndDate(overtime.UserID, overtime.Date)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, errors.New("overtime already exists for today")
	}

	if overtime.Hours > maxHourLimit {
		return nil, errors.New("overtime hours cannot exceed 3 hours per day")
	}
	// Remove cache key if exists
	err = utils.DeleteCacheByPattern(ctx, s.cache, "overtime:*")
	if err != nil {
		return nil, err
	}
	return s.repo.CreateOvertime(ctx, overtime)
}

func (s *overtimeService) UpdateOvertime(ctx context.Context, overtime *models.Overtime) (*models.Overtime, error) {
	cacheKey := utils.BuildKey("overtime", overtime.ID)

	if overtime.Hours > maxHourLimit {
		return nil, errors.New("overtime hours cannot exceed 3 hours per day")
	}

	// Update in DB
	updatedOvertime, err := s.repo.UpdateOvertime(ctx, overtime)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	s.cache.Del(ctx, cacheKey)

	return updatedOvertime, nil
}

func (s *overtimeService) DeleteOvertime(id uint) error {
	ctx := context.Background()
	cacheKey := utils.BuildKey("overtime", id)

	// Delete from DB
	if err := s.repo.DeleteOvertime(id); err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Del(ctx, cacheKey)

	return nil
}
