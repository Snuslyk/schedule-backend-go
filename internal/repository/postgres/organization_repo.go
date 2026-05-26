package postgres

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OrganizationRepository struct {
	db *sqlx.DB
}

func NewOrganizationRepository(db *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `INSERT INTO organizations (name, short_name) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, org.Name, org.ShortName).Scan(&org.ID)
}

func (r *OrganizationRepository) CreateMany(ctx context.Context, orgs []domain.CreateOrganizationDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO organizations (name, short_name) VALUES ($1, $2)`
	for _, o := range orgs {
		if _, err := tx.ExecContext(ctx, query, o.Name, o.ShortName); err != nil {
			return err
		}
	}

	return tx.Commit()
}
