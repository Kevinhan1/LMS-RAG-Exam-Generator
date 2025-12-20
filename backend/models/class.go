package models

type Class struct {
	ID           int64  `json:"id"`
	Jenjang      string `json:"jenjang"`
	Grade        int    `json:"grade"`
	Name         string `json:"name"`
	TimeCreated  int64  `json:"timecreated"`
	TimeModified int64  `json:"timemodified"`
}
