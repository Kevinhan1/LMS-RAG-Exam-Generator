package models

type Question struct {
	ID            int64  `json:"id"`
	MaterialID    int64  `json:"material_id"`
	CreatedBy     int64  `json:"created_by"`
	Content       string `json:"content"`
	Difficulty    string `json:"difficulty"`
	TaxonomyLevel string `json:"taxonomy_level"`
	Status        string `json:"status"`
	TimeCreated   int64  `json:"timecreated"`
	TimeModified  int64  `json:"timemodified"`
}
