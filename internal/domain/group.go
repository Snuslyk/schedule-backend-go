package domain

type Group struct {
	ID   int    `json:"id" db:"id" example:"1"`
	Name string `json:"name" db:"name" example:"ИС-21"`
}

type CreateGroupDTO struct {
	Name string `json:"name" example:"ИС-21"`
}
