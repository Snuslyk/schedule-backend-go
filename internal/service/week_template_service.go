package service

import (
	"context"
	"fmt"

	"github.com/Snuslyk/schedule-backend-go/internal/cache"
	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type WeekTemplateService struct {
	repo      *postgres.WeekTemplateRepository
	groupRepo *postgres.GroupRepository
	cache     *cache.ScheduleCache
}

func NewWeekTemplateService(
	repo *postgres.WeekTemplateRepository,
	groupRepo *postgres.GroupRepository,
	cache *cache.ScheduleCache,
) *WeekTemplateService {
	return &WeekTemplateService{
		repo:      repo,
		groupRepo: groupRepo,
		cache:     cache,
	}
}

func (s *WeekTemplateService) Create(ctx context.Context, groupID int, dto domain.CreateWeekTemplateDTO) (*domain.WeekTemplate, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("There is no group with id %d!", groupID)
	}

	wt, err := s.repo.Create(ctx, groupID, dto)
	if err != nil {
		return nil, err
	}

	s.cache.InvalidateGroup(ctx, groupID, group.Name)
	return wt, nil
}

func (s *WeekTemplateService) CreateMany(ctx context.Context, groupID int, dtos []domain.CreateWeekTemplateDTO) ([]domain.WeekTemplate, error) {
	results := make([]domain.WeekTemplate, 0, len(dtos))
	for _, dto := range dtos {
		created, err := s.Create(ctx, groupID, dto)
		if err != nil {
			return nil, err
		}
		results = append(results, *created)
	}
	return results, nil
}

func (s *WeekTemplateService) DeleteByID(ctx context.Context, id int) (bool, error) {
	return s.repo.DeleteByID(ctx, id)
}
