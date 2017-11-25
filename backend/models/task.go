package models

type Task struct {
	ID uint `gorm:"primary_key" json:"id"`

	ProjectID   uint   `gorm:"index"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
	Description string `json:"Description"`
	Image       string `json:"Image"`
}
