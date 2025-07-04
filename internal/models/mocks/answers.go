package mocks

import (
	"database/sql"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type AnswerModel struct{}

func (m *AnswerModel) InsertAnswers(idsToQuestion models.QuestionJSONSchemaMap, tx *sql.Tx) error {
	return nil
}

func (m *AnswerModel) GetAnswersByQuestionIDs(ids []int) (models.AnswersByQuestion, error) {
	return make(map[int][]models.AnswerPublic), nil
}
