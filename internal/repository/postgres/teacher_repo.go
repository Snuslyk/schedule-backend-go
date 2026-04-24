package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TeacherRepository struct {
	db *sqlx.DB
}

func NewTeacherRepository(db *sqlx.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(ctx context.Context, t *domain.Teacher) error {
	query := `INSERT INTO teachers (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, t.Name).Scan(&t.ID)
}

func (r *TeacherRepository) CreateMany(ctx context.Context, dtos []domain.CreateTeacherDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO teachers (name) VALUES ($1)`
	for _, dto := range dtos {
		if _, err := tx.ExecContext(ctx, query, dto.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TeacherRepository) FindAll(ctx context.Context) ([]domain.Teacher, error) {
	var teachers []domain.Teacher
	query := `SELECT id, name FROM teachers ORDER BY name ASC`
	err := r.db.SelectContext(ctx, &teachers, query)
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r *TeacherRepository) FindByID(ctx context.Context, id int) (*domain.Teacher, error) {
	var t domain.Teacher
	query := `SELECT id, name FROM teachers WHERE id = $1`
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TeacherRepository) DeleteByID(ctx context.Context, id int) (bool, error) {
	query := `DELETE FROM teachers WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}
