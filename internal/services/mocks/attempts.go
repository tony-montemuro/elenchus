package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type AttemptService struct{}

func (s *AttemptService) SaveAttempt(attempt models.AttemptPublic) (models.AttemptPublic, error) {
	return models.AttemptPublic{}, nil
}

func (s *AttemptService) GetAttempt(attemptID, quizID int, profileID *int) (models.AttemptPublic, error) {
	return models.AttemptPublic{}, nil
}
