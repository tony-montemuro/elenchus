package models

import (
	"database/sql"
	"time"
)

type QuestionModelInterface interface {
	InsertQuestions([]QuestionJSONSchema, int, int, *sql.Tx) (QuestionJSONSchemaMap, error)
	GetQuestionsByQuizID(int) ([]QuestionPublic, error)
}

type Question struct {
	ID      int
	Type    QuestionType
	Answers []Answer
	Content string
	Points  uint32
	Created time.Time
	Update  time.Time
	Deleted *time.Time
}

type QuestionPublic struct {
	ID      int
	Type    QuestionTypePublic
	Content string
	Points  int
	Answers []AnswerPublic
}

type QuestionModel struct {
	DB *sql.DB
}

type QuestionJSONSchema struct {
	Content string             `json:"content" jsonschema:"The question that the user is requested to answer"`
	Answers []AnswerJSONSchema `json:"answers" jsonschema:"Between 2 and 4 answers to the question, where exactly 1 is marked as correct"`
}

type QuestionJSONSchemaMap map[int]QuestionJSONSchema

func (q *QuestionPublic) ContainsAnswer(ids []int) bool {
	for _, id := range ids {
		found := false

		for _, answer := range q.Answers {
			if answer.ID == id {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (m *QuestionModel) InsertQuestions(questions []QuestionJSONSchema, quizID int, typeId int, tx *sql.Tx) (QuestionJSONSchemaMap, error) {
	idsToQuestion := make(QuestionJSONSchemaMap)
	stmt, err := tx.Prepare(`INSERT INTO question(quiz_id, type_id, content, points, created, updated)
	VALUES (?, ?, ?, 1, UTC_TIMESTAMP(), UTC_TIMESTAMP())`)
	if err != nil {
		return idsToQuestion, err
	}
	defer stmt.Close()

	for _, question := range questions {
		result, err := stmt.Exec(quizID, typeId, question.Content)
		if err != nil {
			return idsToQuestion, err
		}

		questionID, err := result.LastInsertId()
		if err != nil {
			return idsToQuestion, err
		}

		idsToQuestion[int(questionID)] = question
	}

	return idsToQuestion, nil
}

func (m *QuestionModel) GetQuestionsByQuizID(quizID int) ([]QuestionPublic, error) {
	var questions []QuestionPublic

	stmt := `SELECT qt.name, qt.default_points, q.id, q.content, q.points
	FROM question q
	JOIN question_type qt ON q.type_id = qt.id
	WHERE q.quiz_id = ? AND q.deleted IS NULL
	ORDER BY q.id`

	rows, err := m.DB.Query(stmt, quizID)
	if err != nil {
		return questions, err
	}
	if err = rows.Err(); err != nil {
		return questions, err
	}

	defer rows.Close()

	for rows.Next() {
		var question QuestionPublic
		var questionType QuestionTypePublic

		err = rows.Scan(&questionType.Name, &questionType.DefaultPoints, &question.ID, &question.Content, &question.Points)
		if err != nil {
			return questions, err
		}

		question.Type = questionType
		questions = append(questions, question)
	}

	return questions, nil
}

func (m *QuestionModel) UpdateQuestion(question QuestionPublic, tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE question q
	SET q.content = ?, q.points = ?, q.updated = NOW()
	WHERE q.id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(question.Content, question.Points, question.ID)
	return err
}
