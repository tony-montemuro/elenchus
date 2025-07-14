package services

import (
	"database/sql"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type AttemptServiceInterface interface {
	SaveAttempt(models.AttemptPublic) (int, error)
}

type AttemptService struct {
	DB                         *sql.DB
	AttemptModel               *models.AttemptModel
	MultipleChoiceAttemptModel *models.MultipleChoiceAttemptModel
}

func (s *AttemptService) SaveAttempt(attempt models.AttemptPublic) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := s.AttemptModel.InsertAttempt(attempt, tx)
	if err != nil {
		return 0, err
	}

	if err = s.MultipleChoiceAttemptModel.InsertMultipleChoiceAttempts(id, attempt.QuestionAnswer, tx); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}
