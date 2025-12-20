package models

type Course struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Jenjang      string `json:"jenjang"`
	TimeCreated  int64  `json:"timecreated"`
	TimeModified int64  `json:"timemodified"`
}
