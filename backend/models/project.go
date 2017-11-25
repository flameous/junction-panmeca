package models

type Project struct {
	Model
	PatientID    uint   `gorm:"index"`
	Description  string `json:"Description"`
	RelatedTasks []Task `json:"related_tasks"`
}
