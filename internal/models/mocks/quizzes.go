package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type QuizModel struct{}

func (m *QuizModel) Latest() ([]models.QuizMetadata, error) {
	return []models.QuizMetadata{}, nil
}
