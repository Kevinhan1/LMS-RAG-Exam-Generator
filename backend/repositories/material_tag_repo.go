package repositories

import (
	"context"
	"backendLMS/db"
	"backendLMS/models"
)

func AddMaterialTag(ctx context.Context, mt *models.MaterialTag) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO material_tags(material_id, tag_id) VALUES($1,$2) RETURNING id`,
		mt.MaterialID, mt.TagID,
	).Scan(&mt.ID)
}

func GetMaterialTags(ctx context.Context, materialID int64) ([]models.MaterialTag, error) {
	rows, err := db.Pool.Query(ctx, `SELECT id, material_id, tag_id FROM material_tags WHERE material_id=$1`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.MaterialTag
	for rows.Next() {
		var mt models.MaterialTag
		rows.Scan(&mt.ID, &mt.MaterialID, &mt.TagID)
		tags = append(tags, mt)
	}
	return tags, nil
}