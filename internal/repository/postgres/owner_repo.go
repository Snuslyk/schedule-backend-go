package postgres

import (
	"context"
	"strings"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OwnerRepository struct {
	db *sqlx.DB
}

func NewOwnerRepository(db *sqlx.DB) *OwnerRepository {
	return &OwnerRepository{db: db}
}

func (r *OwnerRepository) FindByName(ctx context.Context, name string, quantity int) ([]domain.Owner, error) {
	if quantity <= 0 {
		quantity = 5
	}
	if quantity > 10 {
		quantity = 10
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return r.findByEmptyName(ctx, quantity)
	}

	return r.findBySearchName(ctx, trimmed, quantity)
}

func (r *OwnerRepository) findByEmptyName(ctx context.Context, quantity int) ([]domain.Owner, error) {
	half := (quantity + 1) / 2
	var groups []domain.Group
	groupQuery := `SELECT id, name FROM groups ORDER BY name ASC LIMIT $1`
	if err := r.db.SelectContext(ctx, &groups, groupQuery, half); err != nil {
		return nil, err
	}

	teacherCount := quantity - len(groups)
	var teachers []domain.Teacher
	if teacherCount > 0 {
		teacherQuery := `SELECT id, name FROM teachers ORDER BY name ASC LIMIT $1`
		if err := r.db.SelectContext(ctx, &teachers, teacherQuery, teacherCount); err != nil {
			return nil, err
		}
	}

	owners := make([]domain.Owner, 0, len(groups)+len(teachers))
	for _, g := range groups {
		owners = append(owners, domain.Owner{
			ID:   g.ID,
			Name: g.Name,
			Type: domain.OwnerTypeGroup,
		})
	}
	for _, t := range teachers {
		owners = append(owners, domain.Owner{
			ID:   t.ID,
			Name: t.Name,
			Type: domain.OwnerTypeTeacher,
		})
	}

	return owners, nil
}

func (r *OwnerRepository) findBySearchName(ctx context.Context, name string, quantity int) ([]domain.Owner, error) {
	pattern := "%" + strings.ToLower(name) + "%"

	var groups []domain.Group
	groupQuery := `SELECT id, name FROM groups WHERE LOWER(name) LIKE $1 ORDER BY name ASC LIMIT $2`
	if err := r.db.SelectContext(ctx, &groups, groupQuery, pattern, quantity); err != nil {
		return nil, err
	}

	owners := make([]domain.Owner, 0, quantity)
	for _, g := range groups {
		owners = append(owners, domain.Owner{
			ID:   g.ID,
			Name: g.Name,
			Type: domain.OwnerTypeGroup,
		})
	}

	remaining := quantity - len(owners)
	if remaining <= 0 {
		return owners, nil
	}

	var teachers []domain.Teacher
	teacherQuery := `SELECT id, name FROM teachers WHERE LOWER(name) LIKE $1 ORDER BY name ASC LIMIT $2`
	if err := r.db.SelectContext(ctx, &teachers, teacherQuery, pattern, remaining); err != nil {
		return nil, err
	}

	for _, t := range teachers {
		owners = append(owners, domain.Owner{
			ID:   t.ID,
			Name: t.Name,
			Type: domain.OwnerTypeTeacher,
		})
	}

	return owners, nil
}
