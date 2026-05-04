package postgres

import (
	"context"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/jmoiron/sqlx"
)

type LessonWithDetails struct {
	domain.Lesson
	TeacherName string `db:"teacher_name"`
	SubjectName string `db:"subject_name"`
	GroupName   string `db:"group_name"`
}

type LessonRepository struct {
	db *sqlx.DB
}

func NewLessonRepository(db *sqlx.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) Create(ctx context.Context, l *domain.Lesson) error {
	query := `
		INSERT INTO lessons (group_id, teacher_id, subject_id, classroom, day_of_week, slot_number, slot_length, week_type, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		l.GroupID, l.TeacherID, l.SubjectID, l.Classroom,
		l.DayOfWeek, l.SlotNumber, l.SlotLength, l.WeekType, l.IsAvailable,
	).Scan(&l.ID)
}

func (r *LessonRepository) CreateMany(ctx context.Context, dtos []domain.CreateLessonDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lessons (group_id, teacher_id, subject_id, classroom, day_of_week, slot_number, slot_length, week_type, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	for _, dto := range dtos {
		isAvailable := true
		if dto.IsAvailable != nil {
			isAvailable = *dto.IsAvailable
		}
		slotLength := dto.SlotLength
		if slotLength <= 0 {
			slotLength = 1
		}
		slotNumber := dto.SlotNumber
		if slotNumber <= 0 {
			slotNumber = 1
		}
		weekType := dto.WeekType
		if weekType == "" {
			weekType = domain.WeekTypeEven
		}

		_, err := tx.ExecContext(ctx, query,
			dto.GroupID, dto.TeacherID, dto.SubjectID, dto.Classroom,
			dto.DayOfWeek, slotNumber, slotLength, weekType, isAvailable,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *LessonRepository) FindByGroupAndDay(ctx context.Context, groupID, dayOfWeek int, weekParity string) ([]LessonWithDetails, error) {
	query := `
		SELECT l.id, l.group_id, l.teacher_id, l.subject_id, l.classroom,
		       l.day_of_week, l.slot_number, l.slot_length, l.week_type, l.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM lessons l
		JOIN teachers t ON l.teacher_id = t.id
		JOIN subjects s ON l.subject_id = s.id
		JOIN groups g ON l.group_id = g.id
		WHERE l.group_id = $1
		  AND l.day_of_week = $2
		  AND (l.week_type = $3 OR l.week_type = 'BOTH' OR l.week_type = 'OTHER')
		ORDER BY l.slot_number ASC
	`
	var items []LessonWithDetails
	err := r.db.SelectContext(ctx, &items, query, groupID, dayOfWeek, weekParity)
	return items, err
}

func (r *LessonRepository) FindByGroupWeek(ctx context.Context, groupID int, weekParity string) ([]LessonWithDetails, error) {
	query := `
		SELECT l.id, l.group_id, l.teacher_id, l.subject_id, l.classroom,
		       l.day_of_week, l.slot_number, l.slot_length, l.week_type, l.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM lessons l
		JOIN teachers t ON l.teacher_id = t.id
		JOIN subjects s ON l.subject_id = s.id
		JOIN groups g ON l.group_id = g.id
		WHERE l.group_id = $1
		  AND (l.week_type = $2 OR l.week_type = 'BOTH' OR l.week_type = 'OTHER')
		ORDER BY l.day_of_week ASC, l.slot_number ASC
	`
	var items []LessonWithDetails
	err := r.db.SelectContext(ctx, &items, query, groupID, weekParity)
	return items, err
}

func (r *LessonRepository) FindByTeacherAndDay(ctx context.Context, teacherID, dayOfWeek int, weekParity string) ([]LessonWithDetails, error) {
	query := `
		SELECT l.id, l.group_id, l.teacher_id, l.subject_id, l.classroom,
		       l.day_of_week, l.slot_number, l.slot_length, l.week_type, l.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM lessons l
		JOIN teachers t ON l.teacher_id = t.id
		JOIN subjects s ON l.subject_id = s.id
		JOIN groups g ON l.group_id = g.id
		WHERE l.teacher_id = $1
		  AND l.day_of_week = $2
		  AND (l.week_type = $3 OR l.week_type = 'BOTH' OR l.week_type = 'OTHER')
		ORDER BY l.slot_number ASC
	`
	var items []LessonWithDetails
	err := r.db.SelectContext(ctx, &items, query, teacherID, dayOfWeek, weekParity)
	return items, err
}

func (r *LessonRepository) FindByTeacherWeek(ctx context.Context, teacherID int, weekParity string) ([]LessonWithDetails, error) {
	query := `
		SELECT l.id, l.group_id, l.teacher_id, l.subject_id, l.classroom,
		       l.day_of_week, l.slot_number, l.slot_length, l.week_type, l.is_available,
		       t.name AS teacher_name, s.name AS subject_name, g.name AS group_name
		FROM lessons l
		JOIN teachers t ON l.teacher_id = t.id
		JOIN subjects s ON l.subject_id = s.id
		JOIN groups g ON l.group_id = g.id
		WHERE l.teacher_id = $1
		  AND (l.week_type = $2 OR l.week_type = 'BOTH' OR l.week_type = 'OTHER')
		ORDER BY l.day_of_week ASC, l.slot_number ASC
	`
	var items []LessonWithDetails
	err := r.db.SelectContext(ctx, &items, query, teacherID, weekParity)
	return items, err
}

func (r *LessonRepository) DeleteByID(ctx context.Context, id int) (bool, error) {
	query := `DELETE FROM lessons WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	return rows > 0, err
}
