package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type QuestionTypeModel struct{}

func (m *QuestionTypeModel) GetMultipleChoiceId() (int, error) {
	return 0, nil
}

func (m *QuestionTypeModel) GetMultipleChoice() (models.QuestionTypePublic, error) {
	return models.QuestionTypePublic{}, nil
}
