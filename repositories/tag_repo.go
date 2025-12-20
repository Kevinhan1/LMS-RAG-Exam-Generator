package repositories

import (
	"context"
	"backendLMS/db"
	"backendLMS/models"
)

func CreateTag(ctx context.Context, t *models.Tag) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO tags (name, description) VALUES ($1,$2) RETURNING id`,
		t.Name, t.Description,
	).Scan(&t.ID)
}

func GetTags(ctx context.Context) ([]models.Tag, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, description FROM tags ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		rows.Scan(&t.ID, &t.Name, &t.Description)
		tags = append(tags, t)
	}
	return tags, nil
}

func UpdateTag(ctx context.Context, t *models.Tag) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE tags SET name=$1, description=$2 WHERE id=$3`,
		t.Name, t.Description, t.ID)
	return err
}

func DeleteTag(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM tags WHERE id=$1`, id)
	return err
}

func GetTagByID(ctx context.Context, id int64) (*models.Tag, error) {
	var t models.Tag
	err := db.Pool.QueryRow(ctx,
		`SELECT id, name, description FROM tags WHERE id=$1`,
		id,
	).Scan(&t.ID, &t.Name, &t.Description)
	return &t, err
}
