package repositories

import (
	"context"
	"backendLMS/db"
	"backendLMS/models"
)

func GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id,name,email,role_id,timecreated FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.RoleID, &u.TimeCreated)
		users = append(users, u)
	}
	return users, nil
}

func UpdateUserRole(ctx context.Context, userID, roleID int64) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE users SET role_id=$1 WHERE id=$2`,
		roleID, userID,
	)
	return err
}

func DeleteUser(ctx context.Context, userID int64) error {
	_, err := db.Pool.Exec(ctx,
		`DELETE FROM users WHERE id=$1`,
		userID,
	)
	return err
}
