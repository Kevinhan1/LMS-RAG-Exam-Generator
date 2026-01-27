package models

type RolePermission struct {
	ID           int64 `json:"id"`
	RoleID       int64 `json:"role_id"`
	PermissionID int64 `json:"permission_id"`
	TimeCreated  int64 `json:"timecreated"`
}
