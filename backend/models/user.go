package models

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	RoleID       int64  `json:"role_id"`
	TimeCreated  int64  `json:"timecreated"`
}