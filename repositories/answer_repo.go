package repositories

import (
	"context"
	"errors"

	"backendLMS/db"
)

var ErrQuestionApproved = errors.New("question already approved, answers cannot be modified")

// ==========================
// helper
// ==========================
func isQuestionApprovedByAnswerID(ctx context.Context, answerID int64) (bool, error) {
	var status string
	err := db.Pool.QueryRow(ctx, `
		SELECT q.status
		FROM questions q
		JOIN answers a ON a.question_id = q.id
		WHERE a.id = $1
	`, answerID).Scan(&status)

	if err != nil {
		return false, err
	}

	return status == "approved", nil
}

func isQuestionApprovedByQuestionID(ctx context.Context, questionID int64) (bool, error) {
	var status string
	err := db.Pool.QueryRow(ctx, `
		SELECT status FROM questions WHERE id = $1
	`, questionID).Scan(&status)

	if err != nil {
		return false, err
	}

	return status == "approved", nil
}

// ==========================
// CREATE
// ==========================
func CreateAnswer(ctx context.Context, questionID int64, input AnswerInput) error {
	approved, err := isQuestionApprovedByQuestionID(ctx, questionID)
	if err != nil {
		return err
	}
	if approved {
		return ErrQuestionApproved
	}

	var count int
	err = db.Pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM answers
		WHERE question_id = $1 AND is_correct = true
	`, questionID).Scan(&count)

	if err != nil {
		return err
	}

	if input.IsCorrect && count > 0 {
		return errors.New("only one correct answer allowed per question")
	}

	_, err = db.Pool.Exec(ctx, `
		INSERT INTO answers (question_id, option_label, option_text, is_correct)
		VALUES ($1,$2,$3,$4)
	`, questionID, input.Label, input.Text, input.IsCorrect)

	return err
}

// ==========================
// UPDATE
// ==========================
func UpdateAnswer(ctx context.Context, answerID int64, input AnswerInput) error {
	approved, err := isQuestionApprovedByAnswerID(ctx, answerID)
	if err != nil {
		return err
	}
	if approved {
		return ErrQuestionApproved
	}

	var questionID int64
	err = db.Pool.QueryRow(ctx, `
		SELECT question_id FROM answers WHERE id=$1
	`, answerID).Scan(&questionID)

	if err != nil {
		return err
	}

	if input.IsCorrect {
		var count int
		err = db.Pool.QueryRow(ctx, `
			SELECT COUNT(*)
			FROM answers
			WHERE question_id=$1 AND is_correct=true AND id<>$2
		`, questionID, answerID).Scan(&count)

		if err != nil {
			return err
		}

		if count > 0 {
			return errors.New("only one correct answer allowed per question")
		}
	}

	_, err = db.Pool.Exec(ctx, `
		UPDATE answers
		SET option_label=$1, option_text=$2, is_correct=$3
		WHERE id=$4
	`, input.Label, input.Text, input.IsCorrect, answerID)

	return err
}

// ==========================
// DELETE
// ==========================
func DeleteAnswer(ctx context.Context, answerID int64) error {
	approved, err := isQuestionApprovedByAnswerID(ctx, answerID)
	if err != nil {
		return err
	}
	if approved {
		return ErrQuestionApproved
	}

	_, err = db.Pool.Exec(ctx, `
		DELETE FROM answers WHERE id=$1
	`, answerID)

	return err
}
