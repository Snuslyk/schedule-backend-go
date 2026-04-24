package domain

type Subject struct {
	ID   int    `json:"id" db:"id" example:"1"`
	Name string `json:"name" db:"name" example:"Математика"`
}

type CreateSubjectDTO struct {
	Name string `json:"name" example:"Математика"`
}
