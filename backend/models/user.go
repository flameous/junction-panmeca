package models

type User interface {
	GetID() uint
	SetProjects([]Project)
}

type Patient struct {
	ID uint `gorm:"primary_key" json:"id"`

	FirstName       string           `json:"FirstName"`
	LastName        string           `json:"LastName"`
	BirthDate       string           `json:"BirthDate"`
	RelatedProjects []Project        `json:"RelatedProjects"`
	ExtraData       PatientExtraData `json:"ExtraData"`
}

func (p *Patient) SetProjects(projects []Project) {
	p.RelatedProjects = projects
}

func (p *Patient) GetID() uint {
	return p.ID
}

type Doctor struct {
	ID uint `gorm:"primary_key" json:"id"`

	FirstName       string          `json:"FirstName"`
	LastName        string          `json:"LastName"`
	BirthDate       string          `json:"BirthDate"`
	RelatedProjects []Project       `json:"RelatedProjects"`
	ExtraData       DoctorExtraData `json:"ExtraData"`
}

func (d *Doctor) SetProjects(projects []Project) {
	d.RelatedProjects = projects
}

func (d *Doctor) GetID() uint {
	return d.ID
}

type PatientExtraData struct {
	Int           int    `json:"Int"`
	String        string `json:"String"`
	IsCoolPatient bool   `json:"IsCoolPatient"`
}

type DoctorExtraData struct {
	Int          int    `json:"Int"`
	String       string `json:"String"`
	IsCoolDoctor bool   `json:"IsCoolDoctor"`
}
