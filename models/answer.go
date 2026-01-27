package models

type Answer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Label      string `json:"label"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}
