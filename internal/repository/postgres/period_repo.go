package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PeriodRepository struct {
	db *sqlx.DB
}

func NewPeriodRepository(db *sqlx.DB) *PeriodRepository {
	return &PeriodRepository{db: db}
}

func (r *PeriodRepository) Create(ctx context.Context, p *domain.AcademicPeriod) error {
	query := `
		INSERT INTO academic_periods (name, start_date, end_date, mode, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, p.Name, p.StartDate, p.EndDate, p.Mode, p.IsActive).Scan(&p.ID)
}

func (r *PeriodRepository) FindActive(ctx context.Context) (*domain.AcademicPeriod, error) {
	var p domain.AcademicPeriod
	query := `
		SELECT id, name, start_date, end_date, mode, is_active
		FROM academic_periods
		WHERE is_active = TRUE
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &p, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *PeriodRepository) FindForDate(ctx context.Context, date time.Time) (*domain.AcademicPeriod, error) {
	var p domain.AcademicPeriod
	query := `
		SELECT id, name, start_date, end_date, mode, is_active
		FROM academic_periods
		WHERE $1 >= start_date AND $1 <= end_date
		ORDER BY id DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &p, query, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
