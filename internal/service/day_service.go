package service

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type DayService struct {
	repo *postgres.DayRepository
}

func NewDayService(repo *postgres.DayRepository) *DayService {
	return &DayService{repo: repo}
}

func (s *DayService) GetByID(ctx context.Context, id int) (*domain.Day, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DayService) DeleteByID(ctx context.Context, id int) (bool, error) {
	return s.repo.DeleteByID(ctx, id)
}
