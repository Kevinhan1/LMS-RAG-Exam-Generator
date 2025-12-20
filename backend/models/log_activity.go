package models

import "time"

type LogActivity struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Action      string    `json:"action"`
	TargetTable string    `json:"target_table"`
	TargetID    int64     `json:"target_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}