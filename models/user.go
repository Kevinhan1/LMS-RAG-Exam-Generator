package models

type User struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	RoleID       int64   `json:"role_id"`
	ClassID      *int64  `json:"class_id,omitempty"`
	NISN         *string `json:"nisn,omitempty"`
	NIP          *string `json:"nip,omitempty"`
	TimeCreated  int64   `json:"timecreated"`
}
