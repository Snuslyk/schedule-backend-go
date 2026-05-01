package domain

type Slot struct {
	ID          int  `json:"id" db:"id" example:"1"`
	SlotNumber  int  `json:"slotNumber" db:"slot_number" example:"1"`
	StartMinute int  `json:"start" db:"start_minute" example:"510"`
	EndMinute   int  `json:"end" db:"end_minute" example:"600"`
	DayOfWeek   *int `json:"dayOfWeek,omitempty" db:"day_of_week" example:"0"`
}

type CreateSlotDTO struct {
	SlotNumber  int  `json:"slotNumber" example:"1"`
	StartMinute int  `json:"start" example:"510"`
	EndMinute   int  `json:"end" example:"600"`
	DayOfWeek   *int `json:"dayOfWeek,omitempty" example:"0"`
}
