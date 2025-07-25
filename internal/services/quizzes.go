package services

import (
	"database/sql"
	"fmt"

	"github.com/tony-montemuro/elenchus/internal/models"
)

type QuizServiceInterface interface {
	UploadQuiz(models.QuizJSONSchema, int) (int, error)
	GetQuizByID(int, *int) (models.QuizPublic, error)
	SaveQuiz(models.QuizPublic, models.QuizPublic) error
	SaveAndPublishQuiz(models.QuizPublic, models.QuizPublic) error
	UnpublishQuizByID(int) error
}

type QuizService struct {
	DB                         *sql.DB
	QuizModel                  *models.QuizModel
	QuestionModel              *models.QuestionModel
	AnswerModel                *models.AnswerModel
	AttemptModel               *models.AttemptModel
	MultipleChoiceAttemptModel *models.MultipleChoiceAttemptModel
	QuestionTypeModel          *models.QuestionTypeModel
}

func (s *QuizService) UploadQuiz(quiz models.QuizJSONSchema, profileID int) (int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	quizID, err := s.QuizModel.InsertQuiz(quiz, profileID, tx)
	if err != nil {
		return 0, err
	}

	typeId, err := s.QuestionTypeModel.GetMultipleChoiceId()
	if err != nil {
		return 0, err
	}

	questionMap, err := s.QuestionModel.InsertQuestions(quiz.Questions, quizID, typeId, tx)
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

func (s *QuizService) SaveQuiz(oldQuiz, newQuiz models.QuizPublic) error {
	return s.saveQuiz(oldQuiz, newQuiz, false)
}

func (s *QuizService) SaveAndPublishQuiz(oldQuiz, newQuiz models.QuizPublic) error {
	return s.saveQuiz(oldQuiz, newQuiz, true)
}

func (s *QuizService) UnpublishQuizByID(id int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.MultipleChoiceAttemptModel.DeleteAttemptsByQuizID(id, tx)
	if err != nil {
		return err
	}

	err = s.AttemptModel.DeleteAttemptsByQuizID(id, tx)
	if err != nil {
		return err
	}

	err = s.QuizModel.UnpublishQuizByID(id, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (s *QuizService) updateAnswers(oldAnswers []models.AnswerPublic, newAnswers []models.AnswerPublic, tx *sql.Tx) error {
	oldAnswerCount, newAnswerCount := len(oldAnswers), len(newAnswers)
	if oldAnswerCount != newAnswerCount {
		return fmt.Errorf("data corruption: old answer count (%d) different from new answer count (%d)", oldAnswerCount, newAnswerCount)
	}

	for i := range oldAnswerCount {
		oldAnswer := oldAnswers[i]
		newAnswer := newAnswers[i]

		if newAnswer.Content != oldAnswer.Content || newAnswer.Correct != oldAnswer.Correct {
			err := s.AnswerModel.UpdateAnswer(newAnswer, tx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *QuizService) updateQuestionsAndAnswers(oldQuestions []models.QuestionPublic, newQuestions []models.QuestionPublic, tx *sql.Tx) error {
	oldQuestionsCount, newQuestionsCount := len(oldQuestions), len(newQuestions)
	if oldQuestionsCount != newQuestionsCount {
		return fmt.Errorf("data corruption: old question count (%d) different from new quiz question count (%d)", oldQuestionsCount, newQuestionsCount)
	}

	for i := range oldQuestionsCount {
		oldQuestion := oldQuestions[i]
		newQuestion := newQuestions[i]

		if newQuestion.Content != oldQuestion.Content || newQuestion.Points != oldQuestion.Points {
			err := s.QuestionModel.UpdateQuestion(newQuestion, tx)
			if err != nil {
				return err
			}
		}

		err := s.updateAnswers(oldQuestion.Answers, newQuestion.Answers, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *QuizService) saveQuiz(oldQuiz, newQuiz models.QuizPublic, isPublish bool) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if newQuiz.Title != oldQuiz.Title || newQuiz.Description != oldQuiz.Description {
		err = s.QuizModel.UpdateQuiz(newQuiz, tx)
		if err != nil {
			return err
		}
	}

	if isPublish {
		err = s.QuizModel.PublishQuizById(newQuiz.ID, tx)
		if err != nil {
			return err
		}
	}

	err = s.updateQuestionsAndAnswers(oldQuiz.Questions, newQuiz.Questions, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
