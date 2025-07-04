package services

import "github.com/tony-montemuro/elenchus/internal/models"

type QuizServiceInterface interface {
	UploadQuiz(models.QuizJSONSchema, int) (int, error)
	GetQuizByID(int, *int) (models.QuizPublic, error)
}

type QuizService struct {
	QuizModel     *models.QuizModel
	QuestionModel *models.QuestionModel
	AnswerModel   *models.AnswerModel
}

func (s *QuizService) UploadQuiz(quiz models.QuizJSONSchema, profileID int) (int, error) {
	tx, err := s.QuizModel.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	quizID, err := s.QuizModel.InsertQuiz(quiz, profileID, tx)
	if err != nil {
		return 0, err
	}

	questionMap, err := s.QuestionModel.InsertQuestions(quiz.Questions, quizID, tx)
	if err != nil {
		return 0, err
	}

	if err = s.AnswerModel.InsertAnswers(questionMap, tx); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return quizID, nil
}

func (s *QuizService) GetQuizByID(id int, profileID *int) (models.QuizPublic, error) {
	quiz, err := s.QuizModel.GetQuizByID(id, profileID)
	if err != nil {
		return quiz, err
	}

	questions, err := s.QuestionModel.GetQuestionsByQuizID(id)
	if err != nil {
		return quiz, err
	}
	quiz.Questions = questions

	questionIDs := make([]int, len(questions))
	for i, question := range questions {
		questionIDs[i] = question.ID
	}

	answersMap, err := s.AnswerModel.GetAnswersByQuestionIDs(questionIDs)
	if err != nil {
		return quiz, nil
	}

	for i, question := range quiz.Questions {
		quiz.Questions[i].Answers = answersMap[question.ID]
	}

	return quiz, err
}
