package service

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type OrganizationService struct {
	repo *postgres.OrganizationRepository
}

func NewOrganizationService(repo *postgres.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) Create(ctx context.Context, dto domain.CreateOrganizationDTO) (*domain.Organization, error) {
	org := &domain.Organization{
		Name:      dto.Name,
		ShortName: dto.ShortName,
	}
	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}
	return org, nil
}

func (s *OrganizationService) CreateMany(ctx context.Context, dtos []domain.CreateOrganizationDTO) (int, error) {
	if err := s.repo.CreateMany(ctx, dtos); err != nil {
		return 0, err
	}
	return len(dtos), nil
}
