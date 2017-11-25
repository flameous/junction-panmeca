package models

type Project struct {
	ID uint `gorm:"primary_key" json:"id"`

	PatientID    uint   `gorm:"index"`
	DoctorID     uint   `gorm:"index"`
	Description  string `json:"Description"`
	RelatedTasks []Task `json:"RelatedTasks"`
}
