package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, g *domain.Group) error {
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
	return r.db.QueryRowContext(ctx, query, g.Name).Scan(&g.ID)
}

func (r *GroupRepository) CreateMany(ctx context.Context, dtos []domain.CreateGroupDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO groups (name) VALUES ($1)`
	for _, dto := range dtos {
		if _, err := tx.ExecContext(ctx, query, dto.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *GroupRepository) FindByID(ctx context.Context, id int) (*domain.Group, error) {
	var g domain.Group
	query := `SELECT id, name FROM groups WHERE id = $1`
	err := r.db.GetContext(ctx, &g, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *GroupRepository) FindByName(ctx context.Context, name string) (*domain.Group, error) {
	var g domain.Group
	query := `SELECT id, name FROM groups WHERE name = $1`
	err := r.db.GetContext(ctx, &g, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *GroupRepository) DeleteByID(ctx context.Context, id int) (bool, error) {
	query := `DELETE FROM groups WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}

func (r *GroupRepository) DeleteByName(ctx context.Context, name string) (bool, error) {
	query := `DELETE FROM groups WHERE name = $1`
	res, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}
