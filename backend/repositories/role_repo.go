package repositories

import (
	"backendLMS/db"
	"backendLMS/models"
	"context"
)

func GetAllRoles(ctx context.Context) ([]models.Role, error) {
	sql := `
	SELECT id, name, description, timecreated, timemodified
	FROM roles
	ORDER BY id
	`

	rows, err := db.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Description,
			&r.TimeCreated,
			&r.TimeModified,
		); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}
