package repositories

import (
	"context"
	"errors"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func CreatePermission(ctx context.Context, p *models.Permission) error {
	now := time.Now().Unix()

	return db.Pool.QueryRow(ctx, `
		INSERT INTO permissions (name, description, timecreated, timemodified)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`, p.Name, p.Description, now, now).Scan(&p.ID)
}

func GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, name, description, timecreated, timemodified
		FROM permissions
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Permission
	for rows.Next() {
		var p models.Permission
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.TimeCreated, &p.TimeModified)
		result = append(result, p)
	}
	return result, nil
}

func GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error) {
	var p models.Permission

	err := db.Pool.QueryRow(ctx, `
		SELECT id, name, description, timecreated, timemodified
		FROM permissions WHERE id=$1
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.TimeCreated, &p.TimeModified)

	if err != nil {
		return nil, errors.New("permission not found")
	}
	return &p, nil
}

func UpdatePermission(ctx context.Context, id int64, p *models.Permission) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE permissions
		SET name=$1, description=$2, timemodified=$3
		WHERE id=$4
	`, p.Name, p.Description, time.Now().Unix(), id)

	return err
}

func DeletePermission(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM permissions WHERE id=$1
	`, id)
	return err
}
