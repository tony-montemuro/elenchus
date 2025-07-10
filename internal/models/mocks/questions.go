package mocks

import (
	"database/sql"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type QuestionModel struct{}

func (m *QuestionModel) InsertQuestions(questions []models.QuestionJSONSchema, quizID int, questionType int, tx *sql.Tx) (models.QuestionJSONSchemaMap, error) {
	return make(models.QuestionJSONSchemaMap), nil
}

func (m *QuestionModel) GetQuestionsByQuizID(quizID int) ([]models.QuestionPublic, error) {
	return []models.QuestionPublic{}, nil
}
