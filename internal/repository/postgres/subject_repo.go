package postgres

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type SubjectRepository struct {
	db *sqlx.DB
}

func NewSubjectRepository(db *sqlx.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) Create(ctx context.Context, s *domain.Subject) error {
	query := `INSERT INTO subjects (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, s.Name).Scan(&s.ID)
}

func (r *SubjectRepository) CreateMany(ctx context.Context, dtos []domain.CreateSubjectDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO subjects (name) VALUES ($1)`
	for _, dto := range dtos {
		if _, err := tx.ExecContext(ctx, query, dto.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SubjectRepository) FindAll(ctx context.Context) ([]domain.Subject, error) {
	var subjects []domain.Subject
	query := `SELECT id, name FROM subjects ORDER BY name ASC`
	err := r.db.SelectContext(ctx, &subjects, query)
	if err != nil {
		return nil, err
	}
	return subjects, nil
}
