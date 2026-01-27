package repositories

import (
	"context"
	"errors"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func CreateMaterial(ctx context.Context, m *models.Material) error {
	now := time.Now().Unix()

	sql := `
	INSERT INTO materials
	    (teacher_id, course_id, chapter_id, title, description, file_url, uploaded_at, timemodified)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id
	`

	return db.Pool.QueryRow(ctx, sql,
		m.TeacherID,
		m.CourseID,
		m.ChapterID,
		m.Title,
		m.Description,
		m.FileURL,
		m.UploadedAt,
		now,
	).Scan(&m.ID)
}

func GetMaterials(ctx context.Context) ([]models.Material, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, teacher_id, course_id, chapter_id,
		       title, description, file_url, uploaded_at, timemodified
		FROM materials
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []models.Material
	for rows.Next() {
		var m models.Material
		_ = rows.Scan(
			&m.ID,
			&m.TeacherID,
			&m.CourseID,
			&m.ChapterID,
			&m.Title,
			&m.Description,
			&m.FileURL,
			&m.UploadedAt,
			&m.TimeModified,
		)
		materials = append(materials, m)
	}

	return materials, nil
}

func GetMaterialByID(ctx context.Context, id int64) (*models.Material, error) {
	var m models.Material

	err := db.Pool.QueryRow(ctx, `
		SELECT id, teacher_id, course_id, chapter_id,
		       title, description, file_url, uploaded_at, timemodified
		FROM materials
		WHERE id = $1
	`, id).Scan(
		&m.ID,
		&m.TeacherID,
		&m.CourseID,
		&m.ChapterID,
		&m.Title,
		&m.Description,
		&m.FileURL,
		&m.UploadedAt,
		&m.TimeModified,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func UpdateMaterial(ctx context.Context, m *models.Material) error {
	now := time.Now().Unix()

	_, err := db.Pool.Exec(ctx, `
		UPDATE materials
		SET title = $1,
		    description = $2,
		    file_url = $3,
		    timemodified = $4
		WHERE id = $5
	`,
		m.Title,
		m.Description,
		m.FileURL,
		now,
		m.ID,
	)

	return err
}

func UpdateMaterialByTeacher(
	ctx context.Context,
	m *models.Material,
	teacherID int64,
) error {
	now := time.Now().Unix()

	cmd, err := db.Pool.Exec(ctx, `
		UPDATE materials
		SET title = $1,
		    description = $2,
		    file_url = $3,
		    timemodified = $4
		WHERE id = $5
		  AND teacher_id = $6
	`,
		m.Title,
		m.Description,
		m.FileURL,
		now,
		m.ID,
		teacherID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("material not found or not owned by teacher")
	}

	return nil
}


func DeleteMaterial(ctx context.Context, id int64) error {
	_, err := db.Pool.Exec(ctx, `
		DELETE FROM materials
		WHERE id = $1
	`, id)

	return err
}

func DeleteMaterialByTeacher(
	ctx context.Context,
	id int64,
	teacherID int64,
) error {
	cmd, err := db.Pool.Exec(ctx, `
		DELETE FROM materials
		WHERE id = $1
		  AND teacher_id = $2
	`, id, teacherID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("material not found or not owned by teacher")
	}

	return nil
}
