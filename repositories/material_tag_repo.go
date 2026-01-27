package repositories

import (
	"context"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func AttachTagToMaterial(ctx context.Context, materialID, tagID int64) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO material_tags (material_id, tag_id, timecreated)
		VALUES ($1,$2,$3)
		ON CONFLICT (material_id, tag_id) DO NOTHING
	`, materialID, tagID, time.Now().Unix())

	return err
}

func GetTagsByMaterial(ctx context.Context, materialID int64) ([]models.Tag, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT t.id, t.name, t.description
		FROM material_tags mt
		JOIN tags t ON mt.tag_id = t.id
		WHERE mt.material_id = $1
	`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
		); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func DetachTagFromMaterial(ctx context.Context, materialID, tagID int64) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM material_tags
		WHERE material_id=$1 AND tag_id=$2
	`, materialID, tagID)

	return err
}
