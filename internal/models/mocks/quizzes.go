package mocks

import (
	"database/sql"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type QuizModel struct{}

func (m *QuizModel) Latest() ([]models.QuizMetadata, error) {
	return []models.QuizMetadata{}, nil
}

func (m *QuizModel) GetQuizByID(id int, profileID *int) (models.QuizPublic, error) {
	return models.QuizPublic{}, nil
}

func (m *QuizModel) InsertQuiz(quiz models.QuizJSONSchema, profileID int, tx *sql.Tx) (int, error) {
	return 0, nil
}

func (m *QuizModel) GetPublishedQuizzesByProfile(profileID *int) ([]models.QuizMetadata, error) {
	return []models.QuizMetadata{}, nil
}

func (m *QuizModel) GetUnpublishedQuizzesByProfile(profileID *int) ([]models.QuizMetadata, error) {
	return []models.QuizMetadata{}, nil
}
