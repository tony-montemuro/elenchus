package models

import (
	"database/sql"
	"fmt"
	"time"
)

type AnswerModelInterface interface {
	InsertAnswers(QuestionJSONSchemaMap, *sql.Tx) error
	GetAnswersByQuestionIDs([]int) (AnswersByQuestion, error)
}

type Answer struct {
	ID      int
	Content string
	Correct bool
	Created time.Time
	Update  time.Time
	Deleted *time.Time
}

type AnswerPublic struct {
	ID      int
	Content string
	Correct bool
}

type AnswersByQuestion map[int][]AnswerPublic

type AnswerModel struct {
	DB *sql.DB
}

type AnswerJSONSchema struct {
	Content string `json:"content" jsonschema:"One of N answers the user has to select from, where N = (2, 3, 4)"`
	Correct bool   `json:"correct" jsonschema:"True if the answer is correct, false otherwise. This flag should be true on exactly one answer per question."`
}

func (m *AnswerModel) InsertAnswers(questionMap QuestionJSONSchemaMap, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT INTO answer(question_id, content, correct, created, updated)
	VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for questionID, question := range questionMap {
		for _, answer := range question.Answers {
			if _, err := stmt.Exec(questionID, answer.Content, answer.Correct); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *AnswerModel) GetAnswersByQuestionIDs(ids []int) (AnswersByQuestion, error) {
	answersMap := make(AnswersByQuestion)
	placeholders, args := buildInClause(ids)
	stmt := fmt.Sprintf(`SELECT a.id, a.content, a.correct, a.question_id
	FROM answer a
	WHERE a.question_id IN (%s)
	ORDER BY a.question_id, a.id`, placeholders)

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return answersMap, err
	}
	if err = rows.Err(); err != nil {
		return answersMap, err
	}

	defer rows.Close()

	for rows.Next() {
		var answer AnswerPublic
		var questionID int

		err = rows.Scan(&answer.ID, &answer.Content, &answer.Correct, &questionID)
		if err != nil {
			return answersMap, err
		}

		if _, ok := answersMap[questionID]; !ok {
			answersMap[questionID] = []AnswerPublic{}
		}
		answersMap[questionID] = append(answersMap[questionID], answer)
	}

	return answersMap, nil
}
