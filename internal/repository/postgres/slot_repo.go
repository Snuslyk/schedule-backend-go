package postgres

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type SlotRepository struct {
	db *sqlx.DB
}

func NewSlotRepository(db *sqlx.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(ctx context.Context, s *domain.Slot) error {
	query := `
		INSERT INTO slots (slot_number, start_minute, end_minute, day_of_week)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, s.SlotNumber, s.StartMinute, s.EndMinute, s.DayOfWeek).Scan(&s.ID)
}

func (r *SlotRepository) CreateMany(ctx context.Context, dtos []domain.CreateSlotDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO slots (slot_number, start_minute, end_minute, day_of_week)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (slot_number, day_of_week) DO UPDATE 
		SET start_minute = EXCLUDED.start_minute, end_minute = EXCLUDED.end_minute
	`
	for _, dto := range dtos {
		if _, err := tx.ExecContext(ctx, query, dto.SlotNumber, dto.StartMinute, dto.EndMinute, dto.DayOfWeek); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SlotRepository) FindAll(ctx context.Context) ([]domain.Slot, error) {
	var slots []domain.Slot
	query := `SELECT id, slot_number, start_minute, end_minute, day_of_week FROM slots ORDER BY slot_number ASC`
	err := r.db.SelectContext(ctx, &slots, query)
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *SlotRepository) FindForDay(ctx context.Context, dayOfWeek int) ([]domain.Slot, error) {
	var slots []domain.Slot
	query := `
		SELECT id, slot_number, start_minute, end_minute, day_of_week
		FROM slots
		WHERE day_of_week = $1
		ORDER BY slot_number ASC
	`
	err := r.db.SelectContext(ctx, &slots, query, dayOfWeek)
	if err != nil {
		return nil, err
	}

	if len(slots) > 0 {
		return slots, nil
	}

	generalQuery := `
		SELECT id, slot_number, start_minute, end_minute, day_of_week
		FROM slots
		WHERE day_of_week IS NULL
		ORDER BY slot_number ASC
	`
	err = r.db.SelectContext(ctx, &slots, generalQuery)
	return slots, err
}
