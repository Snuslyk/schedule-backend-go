package domain

type Teacher struct {
	ID   int    `json:"id" db:"id" example:"1"`
	Name string `json:"name" db:"name" example:"Иванов Иван Иванович"`
}

type CreateTeacherDTO struct {
	Name string `json:"name" example:"Иванов Иван Иванович"`
}
