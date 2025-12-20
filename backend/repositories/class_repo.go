package repositories

import (
	"backendLMS/db"
	"backendLMS/models"
	"context"
	"time"
)

func CreateClass(ctx context.Context, c *models.Class) error {
	sql := `
	INSERT INTO classes (jenjang, grade, name, timecreated, timemodified)
	VALUES ($1,$2,$3,$4,$5) RETURNING id
	`
	return db.Pool.QueryRow(ctx, sql,
		c.Jenjang, c.Grade, c.Name,
		c.TimeCreated, c.TimeModified,
	).Scan(&c.ID)
}

func GetClasses(ctx context.Context) ([]models.Class, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, jenjang, grade, name, timecreated, timemodified FROM classes ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Class
	for rows.Next() {
		var c models.Class
		rows.Scan(&c.ID, &c.Jenjang, &c.Grade, &c.Name, &c.TimeCreated, &c.TimeModified)
		result = append(result, c)
	}
	return result, nil
}

// GetClassByID
func GetClassByID(ctx context.Context, id int64) (*models.Class, error) {
	var c models.Class
	err := db.Pool.QueryRow(ctx,
		`SELECT id, jenjang, grade, name, timecreated, timemodified
		 FROM classes WHERE id=$1`, id,
	).Scan(&c.ID, &c.Jenjang, &c.Grade, &c.Name, &c.TimeCreated, &c.TimeModified)

	if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpdateClass
func UpdateClass(ctx context.Context, c *models.Class) error {
	c.TimeModified = time.Now().Unix()
	_, err := db.Pool.Exec(ctx,
		`UPDATE classes
		 SET jenjang=$1, grade=$2, name=$3, timemodified=$4
		 WHERE id=$5`,
		c.Jenjang, c.Grade, c.Name, c.TimeModified, c.ID,
	)
	return err
}

// DeleteClass
func DeleteClass(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx,
		`DELETE FROM classes WHERE id=$1`, id,
	)
	return err
}
