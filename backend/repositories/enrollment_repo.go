package repositories

import (
	"context"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func Enroll(ctx context.Context, userID, courseID int64) error {
	now := time.Now().Unix()

	_, err := db.Pool.Exec(ctx, `
		INSERT INTO user_courses (user_id, course_id, timecreated, timemodified)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (user_id, course_id) DO NOTHING
	`, userID, courseID, now, now)

	return err
}

func GetUserCourses(ctx context.Context, userID int64) ([]models.Course, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT c.id, c.name, c.jenjang, c.timecreated, c.timemodified
		FROM user_courses uc
		JOIN courses c ON uc.course_id = c.id
		WHERE uc.user_id=$1
		ORDER BY uc.timecreated DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Jenjang,
			&c.TimeCreated,
			&c.TimeModified,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func Unenroll(ctx context.Context, userID, courseID int64) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM user_courses
		WHERE user_id=$1 AND course_id=$2
	`, userID, courseID)

	return err
}

func GetCourseStudents(ctx context.Context, courseID int64) ([]models.User, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT u.id, u.name, u.email
		FROM user_courses uc
		JOIN users u ON uc.user_id = u.id
		WHERE uc.course_id=$1
		ORDER BY uc.timecreated ASC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
		); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}
