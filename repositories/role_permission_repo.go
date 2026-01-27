package repositories

import (
	"context"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func AssignPermissionToRole(ctx context.Context, roleID, permissionID int64) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO role_permissions (role_id, permission_id, timecreated)
		VALUES ($1,$2,$3)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`, roleID, permissionID, time.Now().Unix())

	return err
}

func GetPermissionsByRole(ctx context.Context, roleID int64) ([]models.Permission, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT p.id, p.name, p.description, p.timecreated, p.timemodified
		FROM role_permissions rp
		JOIN permissions p ON rp.permission_id = p.id
		WHERE rp.role_id=$1
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

func RemovePermissionFromRole(ctx context.Context, roleID, permissionID int64) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM role_permissions
		WHERE role_id=$1 AND permission_id=$2
	`, roleID, permissionID)

	return err
}
