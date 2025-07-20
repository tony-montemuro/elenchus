package mocks

import "github.com/tony-montemuro/elenchus/internal/models"

type QuizService struct{}

func (s *QuizService) UploadQuiz(quiz models.QuizJSONSchema, profileID int) (int, error) {
	return 0, nil
}

func (s *QuizService) GetQuizByID(id int, profileID *int) (models.QuizPublic, error) {
	return models.QuizPublic{}, nil
}

func (s *QuizService) SaveAndPublishQuiz(oldQuiz, newQuiz models.QuizPublic) error {
	return nil
}

func (s *QuizService) SaveQuiz(oldQuiz, newQuiz models.QuizPublic) error {
	return nil
}

func (s *QuizService) UnpublishQuizByID(id int) error {
	return nil
}
