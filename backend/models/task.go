package models

type Task struct {
	Model
	ProjectID   uint   `gorm:"index"`
	StartDate   string `json:"StartDate"`
	EndDate     string `json:"EndDate"`
	Description string `json:"Description"`
	Image       string `json:"Image"`
}
