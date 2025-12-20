package repositories

import (
	"backendLMS/db"
	"backendLMS/models"
	"context"
)

func CreateCourse(ctx context.Context, c *models.Course) error {
	sql := `
	INSERT INTO courses (name, jenjang, timecreated, timemodified)
	VALUES ($1,$2,$3,$4) RETURNING id
	`
	return db.Pool.QueryRow(ctx, sql,
		c.Name, c.Jenjang, c.TimeCreated, c.TimeModified,
	).Scan(&c.ID)
}

func GetCourses(ctx context.Context) ([]models.Course, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, jenjang, timecreated, timemodified FROM courses`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Course
	for rows.Next() {
		var c models.Course
		rows.Scan(&c.ID, &c.Name, &c.Jenjang, &c.TimeCreated, &c.TimeModified)
		result = append(result, c)
	}
	return result, nil
}

func GetCourseByID(ctx context.Context, id int64) (*models.Course, error) {
	sql := `
	SELECT id, name, jenjang, timecreated, timemodified
	FROM courses WHERE id=$1
	`
	var c models.Course
	err := db.Pool.QueryRow(ctx, sql, id).
		Scan(&c.ID, &c.Name, &c.Jenjang, &c.TimeCreated, &c.TimeModified)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateCourse(ctx context.Context, c *models.Course) error {
	sql := `
	UPDATE courses
	SET name=$1, jenjang=$2, timemodified=$3
	WHERE id=$4
	`
	_, err := db.Pool.Exec(ctx, sql,
		c.Name,
		c.Jenjang,
		c.TimeModified,
		c.ID,
	)
	return err
}

func DeleteCourse(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return err
}
