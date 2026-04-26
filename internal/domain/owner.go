package domain

type OwnerType string

const (
	OwnerTypeGroup   OwnerType = "group"
	OwnerTypeTeacher OwnerType = "teacher"
)

type Owner struct {
	ID   int       `json:"id" example:"1"`
	Name string    `json:"name" example:"ИС-21"`
	Type OwnerType `json:"type" example:"group"`
}
