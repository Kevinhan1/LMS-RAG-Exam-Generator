package repositories

import (
	"context"
	"backendLMS/db"
	"backendLMS/models"
	"time"
)

func CreateLog(ctx context.Context, log *models.LogActivity) error {
	sql := `
	INSERT INTO log_activity(user_id, action, target_table, target_id, description, created_at)
	VALUES($1,$2,$3,$4,$5,$6) RETURNING id
	`
	log.CreatedAt = time.Now()
	return db.Pool.QueryRow(ctx, sql,
		log.UserID, log.Action, log.TargetTable, log.TargetID, log.Description, log.CreatedAt,
	).Scan(&log.ID)
}
