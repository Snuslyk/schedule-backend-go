package postgres

import (
	"context"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ReplaceWithDetails struct {
	domain.Replace
	TeacherName string `db:"teacher_name"`
	SubjectName string `db:"subject_name"`
	GroupName   string `db:"group_name"`
}

type ReplaceRepository struct {
	db *sqlx.DB
}

func NewReplaceRepository(db *sqlx.DB) *ReplaceRepository {
	return &ReplaceRepository{db: db}
}

func (r *ReplaceRepository) Create(ctx context.Context, rep *domain.Replace) error {
	query := `
		INSERT INTO replaces (date, slot_number, group_id, teacher_id, subject_id, classroom, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		rep.Date, rep.SlotNumber, rep.GroupID, rep.TeacherID, rep.SubjectID, rep.Classroom, rep.IsAvailable,
	).Scan(&rep.ID)
}

func (r *ReplaceRepository) CreateMany(ctx context.Context, reps []domain.Replace) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO replaces (date, slot_number, group_id, teacher_id, subject_id, classroom, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, rep := range reps {
		_, err := tx.ExecContext(ctx, query,
			rep.Date, rep.SlotNumber, rep.GroupID, rep.TeacherID, rep.SubjectID, rep.Classroom, rep.IsAvailable,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReplaceRepository) FindByGroupAndDate(ctx context.Context, groupID int, date time.Time) ([]ReplaceWithDetails, error) {
	query := `
		SELECT r.id, r.date, r.slot_number, r.group_id, r.teacher_id, r.subject_id, r.classroom, r.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM replaces r
		JOIN teachers t ON r.teacher_id = t.id
		JOIN subjects s ON r.subject_id = s.id
		JOIN groups g ON r.group_id = g.id
		WHERE r.group_id = $1 AND r.date = $2::date
		ORDER BY r.slot_number ASC
	`
	var items []ReplaceWithDetails
	err := r.db.SelectContext(ctx, &items, query, groupID, date)
	return items, err
}

func (r *ReplaceRepository) FindByGroupAndDateRange(ctx context.Context, groupID int, start, end time.Time) ([]ReplaceWithDetails, error) {
	query := `
		SELECT r.id, r.date, r.slot_number, r.group_id, r.teacher_id, r.subject_id, r.classroom, r.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM replaces r
		JOIN teachers t ON r.teacher_id = t.id
		JOIN subjects s ON r.subject_id = s.id
		JOIN groups g ON r.group_id = g.id
		WHERE r.group_id = $1 AND r.date >= $2::date AND r.date <= $3::date
		ORDER BY r.date ASC, r.slot_number ASC
	`
	var items []ReplaceWithDetails
	err := r.db.SelectContext(ctx, &items, query, groupID, start, end)
	return items, err
}

func (r *ReplaceRepository) FindByTeacherAndDate(ctx context.Context, teacherID int, date time.Time) ([]ReplaceWithDetails, error) {
	query := `
		SELECT r.id, r.date, r.slot_number, r.group_id, r.teacher_id, r.subject_id, r.classroom, r.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM replaces r
		JOIN teachers t ON r.teacher_id = t.id
		JOIN subjects s ON r.subject_id = s.id
		JOIN groups g ON r.group_id = g.id
		WHERE r.teacher_id = $1 AND r.date = $2::date
		ORDER BY r.slot_number ASC
	`
	var items []ReplaceWithDetails
	err := r.db.SelectContext(ctx, &items, query, teacherID, date)
	return items, err
}

func (r *ReplaceRepository) FindByTeacherAndDateRange(ctx context.Context, teacherID int, start, end time.Time) ([]ReplaceWithDetails, error) {
	query := `
		SELECT r.id, r.date, r.slot_number, r.group_id, r.teacher_id, r.subject_id, r.classroom, r.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM replaces r
		JOIN teachers t ON r.teacher_id = t.id
		JOIN subjects s ON r.subject_id = s.id
		JOIN groups g ON r.group_id = g.id
		WHERE r.teacher_id = $1 AND r.date >= $2::date AND r.date <= $3::date
		ORDER BY r.date ASC, r.slot_number ASC
	`
	var items []ReplaceWithDetails
	err := r.db.SelectContext(ctx, &items, query, teacherID, start, end)
	return items, err
}
