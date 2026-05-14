package domain

import "time"

type Replace struct {
	ID          int       `json:"id" db:"id" example:"1"`
	Date        time.Time `json:"date" db:"date" example:"2026-09-07T00:00:00Z"`
	SlotNumber  int       `json:"slotNumber" db:"slot_number" example:"1"`
	GroupID     int       `json:"groupId" db:"group_id" example:"1"`
	TeacherID   int       `json:"teacherId" db:"teacher_id" example:"1"`
	SubjectID   int       `json:"subjectId" db:"subject_id" example:"1"`
	Classroom   string    `json:"classroom" db:"classroom" example:"404"`
	IsAvailable bool      `json:"isAvailable" db:"is_available" example:"true"`
}

type CreateReplaceDTO struct {
	Date        time.Time `json:"date" example:"2026-09-07T00:00:00Z"`
	SlotNumber  int       `json:"slotNumber" example:"1"`
	Classroom   string    `json:"classroom" example:"404"`
	IsAvailable bool      `json:"isAvailable" example:"true"`
	TeacherID   int       `json:"teacherId" example:"1"`
	SubjectID   int       `json:"subjectId" example:"1"`
	GroupID     int       `json:"groupId,omitempty" example:"1"`
}
