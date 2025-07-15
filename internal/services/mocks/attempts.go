package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type AttemptService struct{}

func (s *AttemptService) SaveAttempt(attempt models.AttemptPublic) (int, error) {
	return 0, nil
}

func (s *AttemptService) GetAttempt(attemptID, quizID int, profileID *int) (models.AttemptPublic, error) {
	return models.AttemptPublic{}, nil
}
