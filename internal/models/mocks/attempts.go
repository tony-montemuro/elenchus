package mocks

import (
	"database/sql"
	"time"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type AttemptModel struct{}

func (m *AttemptModel) InsertAttempt(attempt models.AttemptPublic, profileID int, tx *sql.Tx) (int, error) {
	return 0, nil
}

func (m *AttemptModel) GetAttemptById(id int) (models.AttemptPublic, error) {
	return models.AttemptPublic{}, nil
}

func (m *AttemptModel) GetAttemptDetails(id int) (time.Time, int, error) {
	return time.Now(), 0, nil
}

func (m *AttemptModel) GetAttempts(quizID int, profileID int) ([]models.AttemptMetadata, error) {
	return []models.AttemptMetadata{}, nil
}
