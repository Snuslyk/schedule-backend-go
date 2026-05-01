package domain

import "time"

type AcademicPeriod struct {
	ID        int       `json:"id" db:"id" example:"1"`
	Name      string    `json:"name" db:"name" example:"1 Семестр"`
	StartDate time.Time `json:"startDate" db:"start_date" example:"2026-09-01T00:00:00Z"`
	EndDate   time.Time `json:"endDate" db:"end_date" example:"2026-12-31T00:00:00Z"`
	Mode      string    `json:"mode" db:"mode" example:"PARITY"`
	IsActive  bool      `json:"isActive" db:"is_active" example:"true"`
}

type CreateScheduleDTO struct {
	Start   time.Time `json:"start" example:"2026-09-01T00:00:00Z"`
	End     time.Time `json:"end" example:"2026-12-31T00:00:00Z"`
	Mode    string    `json:"mode" example:"PARITY"`
	GroupID int       `json:"groupId,omitempty" example:"1"`
	Name    string    `json:"name,omitempty" example:"1 Семестр"`
}

type NamedObject struct {
	Name string `json:"name" example:"Математика"`
}

type LessonResponse struct {
	ID          int          `json:"id,omitempty" example:"1"`
	Classroom   string       `json:"classroom" example:"301"`
	SlotNumber  int          `json:"slotNumber" example:"1"`
	SlotLength  int          `json:"slotLength" example:"1"`
	IsAvailable bool         `json:"isAvailable" example:"true"`
	Subject     *NamedObject `json:"subject,omitempty"`
	Teacher     *NamedObject `json:"teacher,omitempty"`
	GroupName   string       `json:"groupName,omitempty" example:"ИС-21"`
}

type SlotResponse struct {
	Start int `json:"start" example:"510"`
	End   int `json:"end" example:"600"`
}

type DayScheduleResponse struct {
	Slots   []SlotResponse   `json:"slots"`
	Lessons []LessonResponse `json:"lessons"`
}

type WeekScheduleResponse struct {
	Type     string                `json:"type" example:"EVEN"`
	WeekType string                `json:"weekType" example:"even"`
	Days     []DayScheduleResponse `json:"days"`
}
