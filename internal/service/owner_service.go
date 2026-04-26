package service

import (
	"context"
	"fmt"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type OwnerService struct {
	ownerRepo *postgres.OwnerRepository
}

func NewOwnerService(ownerRepo *postgres.OwnerRepository) *OwnerService {
	return &OwnerService{ownerRepo: ownerRepo}
}

func (s *OwnerService) FindByName(ctx context.Context, name string, quantity int) ([]domain.Owner, error) {
	if quantity > 10 {
		return nil, fmt.Errorf("quantity must be less than or equal to 10")
	}
	return s.ownerRepo.FindByName(ctx, name, quantity)
}
