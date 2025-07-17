package services

import (
	"database/sql"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type AttemptServiceInterface interface {
	SaveAttempt(models.AttemptPublic) (models.AttemptPublic, error)
	GetAttempt(int, int, *int) (models.AttemptPublic, error)
}

type AttemptService struct {
	DB                         *sql.DB
	AttemptModel               *models.AttemptModel
	MultipleChoiceAttemptModel *models.MultipleChoiceAttemptModel
	QuizService                *QuizService
}

func (s *AttemptService) SaveAttempt(attempt models.AttemptPublic) (models.AttemptPublic, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return attempt, err
	}
	defer tx.Rollback()

	id, err := s.AttemptModel.InsertAttempt(attempt, tx)
	if err != nil {
		return attempt, err
	}
	attempt.ID = &id

	if err = s.MultipleChoiceAttemptModel.InsertMultipleChoiceAttempts(id, attempt.Answers, tx); err != nil {
		return attempt, err
	}

	if err := tx.Commit(); err != nil {
		return attempt, err
	}

	return attempt, nil
}

func (s *AttemptService) GetAttempt(attemptID, quizID int, profileID *int) (models.AttemptPublic, error) {
	var attempt models.AttemptPublic
	quiz, err := s.QuizService.GetQuizByID(quizID, profileID)
	if err != nil {
		return attempt, err
	}

	answers, err := s.MultipleChoiceAttemptModel.GetMultipleChoiceAttempts(attemptID)
	if err != nil {
		return attempt, err
	}

	attempt, err = quiz.Grade(answers)
	if err != nil {
		return attempt, err
	}
	attempt.ID = &attemptID

	created, err := s.AttemptModel.GetAttemptCreatedDate(attemptID)
	if err != nil {
		return attempt, err
	}
	attempt.Created = &created

	return attempt, nil
}
