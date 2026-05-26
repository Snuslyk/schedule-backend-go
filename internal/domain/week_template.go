package domain

type WeekTemplate struct {
	ID       int      `json:"id" db:"id" example:"1"`
	GroupID  int      `json:"groupId" db:"group_id" example:"1"`
	Type     WeekType `json:"type" db:"type" example:"EVEN"`
	Days     []Day    `json:"days"`
}

type Day struct {
	ID             int              `json:"id" db:"id" example:"1"`
	WeekTemplateID int              `json:"weekTemplateId,omitempty" db:"week_template_id" example:"1"`
	DayIndex       int              `json:"dayIndex" db:"day_index" example:"0"`
	Lessons        []LessonResponse `json:"lessons"`
	Slots          []SlotResponse   `json:"slots"`
}

type CreateDaySlotDTO struct {
	Start int `json:"start" example:"480"`
	End   int `json:"end" example:"520"`
}

type CreateDayLessonDTO struct {
	Classroom   string `json:"classroom" example:"405"`
	SlotNumber  int    `json:"slotNumber" example:"1"`
	SlotLength  int    `json:"slotLength" example:"2"`
	IsAvailable bool   `json:"isAvailable" example:"true"`
	TeacherID   int    `json:"teacherId" example:"1"`
	SubjectID   int    `json:"subjectId" example:"1"`
}

type CreateDayDTO struct {
	Lessons []CreateDayLessonDTO `json:"lessons"`
	Slots   []CreateDaySlotDTO   `json:"slots"`
}

type CreateWeekTemplateDTO struct {
	Type WeekType       `json:"type" example:"EVEN"`
	Days []CreateDayDTO `json:"days"`
}
