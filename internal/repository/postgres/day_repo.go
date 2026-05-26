package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type DayRepository struct {
	db *sqlx.DB
}

func NewDayRepository(db *sqlx.DB) *DayRepository {
	return &DayRepository{db: db}
}

func (r *DayRepository) GetByID(ctx context.Context, id int) (*domain.Day, error) {
	var day domain.Day
	err := r.db.GetContext(ctx, &day, `SELECT id, week_template_id, day_index FROM days WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	type slotRow struct {
		StartMinute int `db:"start_minute"`
		EndMinute   int `db:"end_minute"`
	}
	var slotRows []slotRow
	err = r.db.SelectContext(ctx, &slotRows, `SELECT start_minute, end_minute FROM day_slots WHERE day_id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	day.Slots = make([]domain.SlotResponse, len(slotRows))
	for i, sr := range slotRows {
		day.Slots[i] = domain.SlotResponse{
			Start: sr.StartMinute,
			End:   sr.EndMinute,
		}
	}

	type lessonRow struct {
		ID          int    `db:"id"`
		Classroom   string `db:"classroom"`
		SlotNumber  int    `db:"slot_number"`
		SlotLength  int    `db:"slot_length"`
		IsAvailable bool   `db:"is_available"`
		SubjectName string `db:"subject_name"`
		TeacherName string `db:"teacher_name"`
	}
	var lessonRows []lessonRow
	query := `
		SELECT l.id, l.classroom, l.slot_number, l.slot_length, l.is_available,
		       s.name AS subject_name, t.name AS teacher_name
		FROM lessons l
		JOIN subjects s ON l.subject_id = s.id
		JOIN teachers t ON l.teacher_id = t.id
		WHERE l.day_id = $1
		ORDER BY l.slot_number
	`
	err = r.db.SelectContext(ctx, &lessonRows, query, id)
	if err != nil {
		return nil, err
	}

	day.Lessons = make([]domain.LessonResponse, len(lessonRows))
	for i, lr := range lessonRows {
		day.Lessons[i] = domain.LessonResponse{
			ID:          lr.ID,
			Classroom:   lr.Classroom,
			SlotNumber:  lr.SlotNumber,
			SlotLength:  lr.SlotLength,
			IsAvailable: lr.IsAvailable,
			Subject:     &domain.NamedObject{Name: lr.SubjectName},
			Teacher:     &domain.NamedObject{Name: lr.TeacherName},
		}
	}

	return &day, nil
}

func (r *DayRepository) DeleteByID(ctx context.Context, id int) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM days WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
