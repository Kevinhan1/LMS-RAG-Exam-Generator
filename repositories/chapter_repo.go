package repositories

import (
	"backendLMS/db"
	"backendLMS/models"
	"context"
	"time"
)

func CreateChapter(ctx context.Context, c *models.Chapter) error {
	now := time.Now().Unix()
	err := db.Pool.QueryRow(ctx,
		`INSERT INTO chapters(course_id, title, description, timecreated, timemodified)
		 VALUES($1,$2,$3,$4,$5) RETURNING id`,
		c.CourseID, c.Title, c.Description, now, now,
	).Scan(&c.ID)
	if err == nil {
		c.TimeCreated = now
		c.TimeModified = now
	}
	return err
}

func GetChapters(ctx context.Context, courseID int64) ([]models.Chapter, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, course_id, title, description, timecreated, timemodified
		 FROM chapters WHERE course_id=$1 ORDER BY id`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var c models.Chapter
		rows.Scan(&c.ID, &c.CourseID, &c.Title, &c.Description, &c.TimeCreated, &c.TimeModified)
		chapters = append(chapters, c)
	}
	return chapters, nil
}

func UpdateChapter(ctx context.Context, c *models.Chapter) error {
	now := time.Now().Unix()
	_, err := db.Pool.Exec(ctx,
		`UPDATE chapters SET title=$1, description=$2, timemodified=$3 WHERE id=$4`,
		c.Title, c.Description, now, c.ID)
	if err == nil {
		c.TimeModified = now
	}
	return err
}

func DeleteChapter(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM chapters WHERE id=$1`, id)
	return err
}

func GetChapterByID(ctx context.Context, id int64) (*models.Chapter, error) {
	row := db.Pool.QueryRow(ctx,
		`SELECT id, course_id, title, description, order_no, timecreated, timemodified
		 FROM chapters WHERE id=$1`, id)

	var c models.Chapter
	err := row.Scan(
		&c.ID,
		&c.CourseID,
		&c.Title,
		&c.Description,
		&c.OrderNo,
		&c.TimeCreated,
		&c.TimeModified,
	)
	return &c, err
}
