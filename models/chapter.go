package models

type Chapter struct {
	ID           int64  `json:"id"`
	CourseID     int64  `json:"course_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	OrderNo      int    `json:"order_no"`
	TimeCreated  int64  `json:"timecreated"`
	TimeModified int64  `json:"timemodified"`
}