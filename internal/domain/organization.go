package domain

type Organization struct {
	ID        int    `json:"id" db:"id" example:"1"`
	Name      string `json:"name" db:"name" example:"Лицей №1"`
	ShortName string `json:"shortName" db:"short_name" example:"Л1"`
}

type CreateOrganizationDTO struct {
	Name      string `json:"name" example:"Лицей №1"`
	ShortName string `json:"shortName" example:"Л1"`
}
