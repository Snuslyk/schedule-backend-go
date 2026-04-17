package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, name, email, password, roles, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.Password, pq.Array(u.Roles), u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, name, email, password, roles, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &u, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, name, email, password, roles, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id string, role string) error {
	query := `UPDATE users SET roles = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, pq.Array([]string{role}), id)
	return err
}
