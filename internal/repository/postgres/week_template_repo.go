package postgres

import (
	"context"
	"fmt"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type WeekTemplateRepository struct {
	db *sqlx.DB
}

func NewWeekTemplateRepository(db *sqlx.DB) *WeekTemplateRepository {
	return &WeekTemplateRepository{db: db}
}

func (r *WeekTemplateRepository) Create(ctx context.Context, groupID int, dto domain.CreateWeekTemplateDTO) (*domain.WeekTemplate, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var weekTemplateID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO week_templates (group_id, type) VALUES ($1, $2) RETURNING id`,
		groupID, dto.Type,
	).Scan(&weekTemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to create week template: %w", err)
	}

	createdDays := make([]domain.Day, len(dto.Days))

	for dayIdx, d := range dto.Days {
		var dayID int
		err = tx.QueryRowContext(ctx,
			`INSERT INTO days (week_template_id, day_index) VALUES ($1, $2) RETURNING id`,
			weekTemplateID, dayIdx,
		).Scan(&dayID)
		if err != nil {
			return nil, fmt.Errorf("failed to create day: %w", err)
		}

		daySlots := make([]domain.SlotResponse, len(d.Slots))
		for sIdx, s := range d.Slots {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO day_slots (day_id, start_minute, end_minute) VALUES ($1, $2, $3)`,
				dayID, s.Start, s.End,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create day slot: %w", err)
			}
			daySlots[sIdx] = domain.SlotResponse{
				Start: s.Start,
				End:   s.End,
			}
		}

		dayLessons := make([]domain.LessonResponse, len(d.Lessons))
		for lIdx, l := range d.Lessons {
			var lessonID int
			err = tx.QueryRowContext(ctx,
				`INSERT INTO lessons (group_id, teacher_id, subject_id, classroom, day_of_week, slot_number, slot_length, week_type, is_available, day_id)
				 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				 RETURNING id`,
				groupID, l.TeacherID, l.SubjectID, l.Classroom, dayIdx, l.SlotNumber, l.SlotLength, dto.Type, l.IsAvailable, dayID,
			).Scan(&lessonID)
			if err != nil {
				return nil, fmt.Errorf("failed to create lesson: %w", err)
			}

			var subjectName string
			_ = tx.QueryRowContext(ctx, `SELECT name FROM subjects WHERE id = $1`, l.SubjectID).Scan(&subjectName)
			var teacherName string
			_ = tx.QueryRowContext(ctx, `SELECT name FROM teachers WHERE id = $1`, l.TeacherID).Scan(&teacherName)

			dayLessons[lIdx] = domain.LessonResponse{
				ID:          lessonID,
				Classroom:   l.Classroom,
				SlotNumber:  l.SlotNumber,
				SlotLength:  l.SlotLength,
				IsAvailable: l.IsAvailable,
				Subject:     &domain.NamedObject{Name: subjectName},
				Teacher:     &domain.NamedObject{Name: teacherName},
			}
		}

		createdDays[dayIdx] = domain.Day{
			ID:             dayID,
			WeekTemplateID: weekTemplateID,
			DayIndex:       dayIdx,
			Lessons:        dayLessons,
			Slots:          daySlots,
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &domain.WeekTemplate{
		ID:      weekTemplateID,
		GroupID: groupID,
		Type:    dto.Type,
		Days:    createdDays,
	}, nil
}

func (r *WeekTemplateRepository) DeleteByID(ctx context.Context, id int) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM week_templates WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
