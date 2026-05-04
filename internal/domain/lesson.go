package domain

type WeekType string

const (
	WeekTypeEven  WeekType = "EVEN"
	WeekTypeOdd   WeekType = "ODD"
	WeekTypeBoth  WeekType = "BOTH"
	WeekTypeOther WeekType = "OTHER"
)

type Lesson struct {
	ID          int      `json:"id" db:"id" example:"1"`
	GroupID     int      `json:"groupId" db:"group_id" example:"1"`
	TeacherID   int      `json:"teacherId" db:"teacher_id" example:"1"`
	SubjectID   int      `json:"subjectId" db:"subject_id" example:"1"`
	Classroom   string   `json:"classroom" db:"classroom" example:"301"`
	DayOfWeek   int      `json:"dayOfWeek" db:"day_of_week" example:"0"`
	SlotNumber  int      `json:"slotNumber" db:"slot_number" example:"1"`
	SlotLength  int      `json:"slotLength" db:"slot_length" example:"1"`
	WeekType    WeekType `json:"weekType" db:"week_type" example:"BOTH"`
	IsAvailable bool     `json:"isAvailable" db:"is_available" example:"true"`
}

type CreateLessonDTO struct {
	GroupID     int      `json:"groupId" example:"1"`
	TeacherID   int      `json:"teacherId" example:"1"`
	SubjectID   int      `json:"subjectId" example:"1"`
	Classroom   string   `json:"classroom" example:"301"`
	DayOfWeek   int      `json:"dayOfWeek" example:"0"`
	SlotNumber  int      `json:"slotNumber" example:"1"`
	SlotLength  int      `json:"slotLength" example:"1"`
	WeekType    WeekType `json:"weekType" example:"BOTH"`
	IsAvailable *bool    `json:"isAvailable,omitempty" example:"true"`
}
