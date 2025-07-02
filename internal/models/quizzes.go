package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type QuizModelInterface interface {
	Latest() ([]QuizMetadata, error)
	GetQuizByID(int) (QuizPublic, error)
}

type Quiz struct {
	ID          int
	Profile     Profile
	Title       string
	Description string
	Questions   []Question
	Published   *time.Time
	Unpublished *time.Time
	Created     time.Time
	Update      time.Time
	Deleted     *time.Time
}

type QuizMetadata struct {
	ID            int
	Profile       ProfilePublic
	Title         string
	Description   string
	QuestionCount int
	Published     time.Time
}

type QuizPublic struct {
	ID          int
	Profile     ProfilePublic
	Title       string
	Description string
	Questions   []QuestionPublic
	Published   time.Time
}

type QuizModel struct {
	DB *sql.DB
}

func (m *QuizModel) Latest() ([]QuizMetadata, error) {
	var quizzes []QuizMetadata

	stmt := `SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, (SELECT count(id) FROM question WHERE quiz_id = q.id) AS question_count, q.published
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id
	WHERE q.published IS NOT NULL AND (q.unpublished IS NULL OR q.published > q.unpublished) AND q.deleted IS NULL
	ORDER BY q.published DESC
	`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return quizzes, nil
	}

	defer rows.Close()

	for rows.Next() {
		var p ProfilePublic
		var q QuizMetadata

		err = rows.Scan(&q.ID, &p.ID, &p.FirstName, &p.LastName, &p.Deleted, &q.Title, &q.Description, &q.QuestionCount, &q.Published)
		if err != nil {
			return []QuizMetadata{}, err
		}

		q.Profile = p

		quizzes = append(quizzes, q)
	}

	return quizzes, nil
}

func (m *QuizModel) GetQuizByID(id int) (QuizPublic, error) {
	quiz, err := m.getQuizByID(id)
	if err != nil {
		return quiz, err
	}

	questions, err := m.getQuestionsByQuizID(id)
	if err != nil {
		return quiz, err
	}
	quiz.Questions = questions

	var questionIDs []int
	for _, question := range questions {
		questionIDs = append(questionIDs, question.ID)
	}

	answersByQuestionID, err := m.getAnswersByQuestionIDs(questionIDs)
	if err != nil {
		return quiz, err
	}

	for i, question := range quiz.Questions {
		quiz.Questions[i].Answers = answersByQuestionID[question.ID]
	}

	return quiz, err
}

func (m *QuizModel) getQuizByID(id int) (QuizPublic, error) {
	var quiz QuizPublic
	var profile ProfilePublic

	stmt := `SELECT q.id, p.id, p.first_name, p.last_name, p.deleted, q.title, q.description, q.published
	FROM quiz q
	JOIN profile p ON q.profile_id = p.id 
	WHERE q.id = ? AND q.published IS NOT NULL AND (q.unpublished IS NULL OR q.published > q.unpublished) AND q.deleted IS NULL`

	err := m.DB.QueryRow(stmt, id).Scan(&quiz.ID, &profile.ID, &profile.FirstName, &profile.LastName, &profile.Deleted, &quiz.Title, &quiz.Description, &quiz.Published)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return QuizPublic{}, ErrNoRecord
		} else {
			return QuizPublic{}, err
		}
	}

	quiz.Profile = profile
	return quiz, nil
}

func (m *QuizModel) getQuestionsByQuizID(id int) ([]QuestionPublic, error) {
	var questions []QuestionPublic

	stmt := `SELECT qt.name, qt.default_points, q.id, q.content, q.points
	FROM question q
	JOIN question_type qt ON q.type_id = qt.id
	WHERE q.quiz_id = ? AND q.deleted IS NULL
	ORDER BY q.id`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return questions, err
	}
	if err = rows.Err(); err != nil {
		return questions, err
	}

	defer rows.Close()

	for rows.Next() {
		var question QuestionPublic
		var questionType QuestionTypePublic

		err = rows.Scan(&questionType.Name, &questionType.DefaultPoints, &question.ID, &question.Content, &question.Points)
		if err != nil {
			return questions, err
		}

		question.Type = questionType
		questions = append(questions, question)
	}

	return questions, nil
}

func (m *QuizModel) getAnswersByQuestionIDs(ids []int) (map[int][]AnswerPublic, error) {
	answersByQuestionID := make(map[int][]AnswerPublic)
	placeholders, args := buildInClause(ids)
	stmt := fmt.Sprintf(`SELECT a.content, a.correct, a.question_id
	FROM answer a
	WHERE a.question_id IN (%s)
	ORDER BY a.question_id, a.id`, placeholders)

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return answersByQuestionID, err
	}
	if err = rows.Err(); err != nil {
		return answersByQuestionID, err
	}

	defer rows.Close()

	for rows.Next() {
		var answer AnswerPublic
		var questionID int

		err = rows.Scan(&answer.Content, &answer.Correct, &questionID)
		if err != nil {
			return answersByQuestionID, err
		}

		answersByQuestionID[questionID] = append(answersByQuestionID[questionID], answer)
	}

	return answersByQuestionID, nil
}
