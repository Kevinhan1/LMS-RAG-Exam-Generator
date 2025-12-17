package repositories

import (
	"context"
	"backendLMS/db"
	"backendLMS/models"
)

func CreateUser(ctx context.Context, u *models.User) error {
	sql := `
	INSERT INTO users (name,email,password_hash,role_id,timecreated)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`
	return db.Pool.QueryRow(
		ctx, sql,
		u.Name, u.Email, u.PasswordHash, u.RoleID, u.TimeCreated,
	).Scan(&u.ID)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `
	SELECT id,name,email,password_hash,role_id,timecreated
	FROM users WHERE email=$1
	`
	var u models.User
	err := db.Pool.QueryRow(ctx, sql, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.RoleID, &u.TimeCreated)
	return &u, err
}