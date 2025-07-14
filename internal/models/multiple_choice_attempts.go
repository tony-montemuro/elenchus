package models

import "database/sql"

type MultipleChoiceAttemptModel struct {
	DB *sql.DB
}

func (m *MultipleChoiceAttemptModel) InsertMultipleChoiceAttempts(attemptID int, answers QuestionAnswer, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT INTO multiple_choice_attempt (attempt_id, question_id, answer_id, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`)
	if err != nil {
		return err
	}

	for questionID, answerID := range answers {
		_, err := stmt.Exec(attemptID, questionID, answerID)
		if err != nil {
			return err
		}
	}

	return nil
}
