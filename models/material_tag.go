package models

type MaterialTag struct {
	ID         int64 `json:"id"`
	MaterialID int64 `json:"material_id"`
	TagID      int64 `json:"tag_id"`
	TimeCreated int64 `json:"timecreated"`
}
