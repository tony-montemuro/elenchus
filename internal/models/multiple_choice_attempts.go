package models

import (
	"database/sql"
)

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

func (m *MultipleChoiceAttemptModel) GetMultipleChoiceAttempts(attemptID int) (QuestionAnswer, error) {
	answers := make(QuestionAnswer)
	stmt := `SELECT mca.question_id, mca.answer_id
	FROM multiple_choice_attempt mca
	WHERE mca.attempt_id = ?
	ORDER BY mca.question_id`

	rows, err := m.DB.Query(stmt, attemptID)
	if err != nil {
		return answers, err
	}
	if err = rows.Err(); err != nil {
		return answers, err
	}

	defer rows.Close()

	for rows.Next() {
		var questionID int
		var answerID int

		err = rows.Scan(&questionID, &answerID)
		if err != nil {
			return answers, err
		}

		answers[questionID] = answerID
	}

	return answers, nil

}

func (m *MultipleChoiceAttemptModel) DeleteAttemptsByQuizID(quizID int, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE multiple_choice_attempt mca
	JOIN attempt a ON mca.attempt_id = a.id
	SET mca.deleted = UTC_TIMESTAMP()
	WHERE a.quiz_id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(quizID)
	return err
}
