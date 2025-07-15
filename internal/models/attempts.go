package models

import (
	"database/sql"
	"errors"
	"time"
)

type QuestionAnswer map[int]int

type AttemptPublic struct {
	ID           *int
	Quiz         QuizPublic
	PointsEarned int
	Answers      QuestionAnswer
	Created      *time.Time
}

type AttemptModel struct {
	DB *sql.DB
}

func (m *AttemptModel) InsertAttempt(attempt AttemptPublic, tx *sql.Tx) (int, error) {
	stmt, err := tx.Prepare(`INSERT INTO attempt (profile_id, quiz_id, points_earned, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(attempt.Quiz.Profile.ID, attempt.Quiz.ID, attempt.PointsEarned)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (m *AttemptModel) GetAttemptById(id int) (AttemptPublic, error) {
	var attempt AttemptPublic

	stmt := `SELECT id, points_earned, created
	FROM attempt
	WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&attempt.ID, &attempt.PointsEarned, &attempt.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attempt, ErrNoRecord
		}

		return attempt, err
	}

	return attempt, nil
}
