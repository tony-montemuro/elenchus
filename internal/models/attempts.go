package models

import (
	"database/sql"
	"errors"
	"time"
)

type AttemptModelInterface interface {
	InsertAttempt(AttemptPublic, *sql.Tx) (int, error)
	GetAttemptById(int) (AttemptPublic, error)
	GetAttemptCreatedDate(int) (time.Time, error)
	GetAttempts(int, int) ([]AttemptMetadata, error)
}

type QuestionAnswer map[int]int

type AttemptPublic struct {
	ID           *int
	Quiz         QuizPublic
	PointsEarned int
	Answers      QuestionAnswer
	Created      *time.Time
}

type AttemptMetadata struct {
	ID           int
	Sequence     int
	PointsEarned int
	Created      time.Time
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
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return attempt, ErrNoRecord
	}

	return attempt, err
}

func (m *AttemptModel) GetAttemptCreatedDate(id int) (time.Time, error) {
	var created time.Time

	stmt := `SELECT created
	FROM attempt
	WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&created)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return created, ErrNoRecord
	}

	return created, err
}

func (m *AttemptModel) GetAttempts(quizID, profileID int) ([]AttemptMetadata, error) {
	attempts := []AttemptMetadata{}

	stmt := `SELECT a.id, a.points_earned, a.created, ROW_NUMBER() OVER (ORDER BY a.created ASC) AS sequence 
	FROM attempt a
	WHERE a.quiz_id = ? AND a.profile_id = ?
	ORDER BY a.created DESC`

	rows, err := m.DB.Query(stmt, quizID, profileID)
	if err != nil {
		return attempts, err
	}

	defer rows.Close()

	for rows.Next() {
		var a AttemptMetadata

		err := rows.Scan(&a.ID, &a.PointsEarned, &a.Created, &a.Sequence)
		if err != nil {
			return attempts, err
		}

		attempts = append(attempts, a)
	}

	return attempts, nil
}
