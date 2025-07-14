package models

import "database/sql"

type QuestionAnswer map[int]int

type AttemptPublic struct {
	Quiz           QuizPublic
	PointsEarned   int
	QuestionAnswer QuestionAnswer
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
