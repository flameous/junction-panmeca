package models

type Model struct {
	ID uint `gorm:"primary_key" json:"id"`
}

type Patient struct {
	Model
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	BirthDate       string    `json:"birth_date"`
	RelatedProjects []Project `json:"related_projects"`
}
