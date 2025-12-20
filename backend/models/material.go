package models

import "time"

type Material struct {
	ID          int64     `json:"id"`
	TeacherID   int64     `json:"teacher_id"`
	CourseID    int64     `json:"course_id"`
	ChapterID   int64     `json:"chapter_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileURL     string    `json:"file_url"`
	UploadedAt  time.Time `json:"uploaded_at"`
	TimeModified int64    `json:"timemodified"`
}