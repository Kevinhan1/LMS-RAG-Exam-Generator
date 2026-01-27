package repositories

import (
	"context"
	"errors"
	"time"

	"backendLMS/db"
	"backendLMS/models"
)

func CreateQuestionWithAnswers(
	ctx context.Context,
	materialID, teacherID int64,
	content, difficulty, taxonomy string,
	answers []AnswerInput,
) error {

	if err := validateAnswers(answers); err != nil {
		return err
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now().Unix()
	var questionID int64

	err = tx.QueryRow(ctx, `
		INSERT INTO questions
		(material_id, created_by, content, difficulty, taxonomy_level, status, timecreated, timemodified)
		VALUES ($1,$2,$3,$4,$5,'draft',$6,$6)
		RETURNING id
	`, materialID, teacherID, content, difficulty, taxonomy, now).
		Scan(&questionID)

	if err != nil {
		return err
	}

	for _, a := range answers {
		_, err := tx.Exec(ctx, `
			INSERT INTO answers (question_id, option_label, option_text, is_correct)
			VALUES ($1,$2,$3,$4)
		`, questionID, a.Label, a.Text, a.IsCorrect)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func GetQuestions(ctx context.Context, userID, roleID int64) ([]models.Question, error) {
	var query string
	var args []interface{}

	if roleID == 1 { // ADMIN
		query = `
			SELECT id, material_id, created_by, content,
			       difficulty, taxonomy_level, status,
			       timecreated, timemodified
			FROM questions
			ORDER BY timecreated DESC
		`
	} else { // TEACHER
		query = `
			SELECT id, material_id, created_by, content,
			       difficulty, taxonomy_level, status,
			       timecreated, timemodified
			FROM questions
			WHERE created_by = $1
			ORDER BY timecreated DESC
		`
		args = append(args, userID)
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(
			&q.ID,
			&q.MaterialID,
			&q.CreatedBy,
			&q.Content,
			&q.Difficulty,
			&q.TaxonomyLevel,
			&q.Status,
			&q.TimeCreated,
			&q.TimeModified,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, nil
}

func GetQuestionByID(ctx context.Context, id, userID, roleID int64) (*models.Question, []models.Answer, error) {
	var q models.Question
	var query string
	var args []interface{}

	if roleID == 1 { // ADMIN
		query = `
			SELECT id, material_id, created_by, content,
			       difficulty, taxonomy_level, status,
			       timecreated, timemodified
			FROM questions
			WHERE id=$1
		`
		args = append(args, id)
	} else { // TEACHER
		query = `
			SELECT id, material_id, created_by, content,
			       difficulty, taxonomy_level, status,
			       timecreated, timemodified
			FROM questions
			WHERE id=$1 AND created_by=$2
		`
		args = append(args, id, userID)
	}

	err := db.Pool.QueryRow(ctx, query, args...).Scan(
		&q.ID,
		&q.MaterialID,
		&q.CreatedBy,
		&q.Content,
		&q.Difficulty,
		&q.TaxonomyLevel,
		&q.Status,
		&q.TimeCreated,
		&q.TimeModified,
	)

	if err != nil {
		return nil, nil, err
	}

	rows, err := db.Pool.Query(ctx, `
		SELECT id, question_id, option_label, option_text, is_correct
		FROM answers
		WHERE question_id=$1
		ORDER BY option_label
	`, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var a models.Answer
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.Label, &a.Text, &a.IsCorrect); err != nil {
			return nil, nil, err
		}
		answers = append(answers, a)
	}

	return &q, answers, nil
}

func UpdateQuestion(ctx context.Context, qID, userID, roleID int64,
	content, difficulty, taxonomy string,
	answers []AnswerInput,
) error {

	if err := validateAnswers(answers); err != nil {
		return err
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var query string
	var args []interface{}

	if roleID == 1 { // ADMIN: Can edit everything, no status check usually needed but lets keep logic simple
		query = `
			UPDATE questions
			SET content=$1, difficulty=$2, taxonomy_level=$3, timemodified=$4
			WHERE id=$5
		`
		args = append(args, content, difficulty, taxonomy, time.Now().Unix(), qID)
	} else { // TEACHER: Only own draft questions
		query = `
			UPDATE questions
			SET content=$1, difficulty=$2, taxonomy_level=$3, timemodified=$4
			WHERE id=$5 AND created_by=$6 AND status='draft'
		`
		args = append(args, content, difficulty, taxonomy, time.Now().Unix(), qID, userID)
	}

	res, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("question not found or not editable")
	}

	_, err = tx.Exec(ctx, `DELETE FROM answers WHERE question_id=$1`, qID)
	if err != nil {
		return err
	}

	for _, a := range answers {
		_, err := tx.Exec(ctx, `
			INSERT INTO answers (question_id, option_label, option_text, is_correct)
			VALUES ($1,$2,$3,$4)
		`, qID, a.Label, a.Text, a.IsCorrect)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func DeleteQuestion(ctx context.Context, qID, userID, roleID int64) error {
	var query string
	var args []interface{}

	if roleID == 1 { // ADMIN
		query = `DELETE FROM questions WHERE id=$1`
		args = append(args, qID)
	} else { // TEACHER
		query = `DELETE FROM questions WHERE id=$1 AND created_by=$2 AND status='draft'`
		args = append(args, qID, userID)
	}

	res, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("question not found or not deletable")
	}

	return nil
}

func UpdateQuestionStatus(ctx context.Context, qID, userID, roleID int64, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	var checkQuery string
	var checkArgs []interface{}

	if roleID == 1 {
		checkQuery = `SELECT EXISTS(SELECT 1 FROM questions WHERE id=$1)`
		checkArgs = append(checkArgs, qID)
	} else {
		checkQuery = `SELECT EXISTS(SELECT 1 FROM questions WHERE id=$1 AND created_by=$2)`
		checkArgs = append(checkArgs, qID, userID)
	}

	err = tx.QueryRow(ctx, checkQuery, checkArgs...).Scan(&exists)
	if err != nil || !exists {
		return errors.New("question not found")
	}

	_, err = tx.Exec(ctx, `
		UPDATE questions
		SET status=$1, timemodified=$2
		WHERE id=$3
	`, status, time.Now().Unix(), qID)

	if err != nil {
		return err
	}

	// AUTO INSERT TO QUESTION_BANK
	if status == "approved" {
		_, err = tx.Exec(ctx, `
			INSERT INTO question_bank (question_id, approved_by, approved_at)
			VALUES ($1,$2,now())
			ON CONFLICT (question_id) DO NOTHING
		`, qID, userID) // userID here is whoever doing action (approver)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func validateAnswers(answers []AnswerInput) error {
	if len(answers) < 2 {
		return errors.New("minimal 2 jawaban diperlukan")
	}

	correct := 0
	for _, a := range answers {
		if a.IsCorrect {
			correct++
		}
	}

	if correct != 1 {
		return errors.New("harus tepat 1 jawaban benar")
	}

	return nil
}
