package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/tony-montemuro/elenchus/internal/models"
	"github.com/tony-montemuro/elenchus/internal/validator"
)

type signupForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Password2 string
	validator.Validator
}

func (f signupForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["firstName"] = f.FirstName
	vals["lastName"] = f.LastName
	vals["email"] = f.Email
	vals["password"] = f.Password

	return vals
}

type loginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (f loginForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["email"] = f.Email
	vals["password"] = f.Password

	return vals
}

type createForm struct {
	Notes string
	validator.Validator
}

func (f createForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["notes"] = f.Notes

	return vals
}

type profileForm struct {
	FirstName string
	LastName  string
	validator.Validator
}

func (f profileForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["firstName"] = f.FirstName
	vals["lastName"] = f.LastName

	return vals
}

type questionEdit struct {
	Content string
	Correct int
	Points  int
	Answers []answerEdit
}

type answerEdit struct {
	ID      int
	Content string
}

type questionEditMap map[int]questionEdit
type answerEditMap map[int]answerEdit
type answersEditMap map[int][]answerEdit

type editForm struct {
	Title       string
	Description string
	Questions   questionEditMap
	validator.Validator
}

func newEditForm(postForm url.Values) (editForm, error) {
	editForm := editForm{}
	err := editForm.parseRequest(postForm)
	return editForm, err
}

func (f editForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["title"] = f.Title
	vals["description"] = f.Description

	return vals
}

func (f *editForm) parseRequest(postForm url.Values) error {
	f.Questions = make(questionEditMap)

	for key, vals := range postForm {
		if len(vals) == 0 || key == "csrf_token" || strings.HasPrefix(key, "answer[") {
			continue
		}

		val := vals[0]
		if strings.HasPrefix(key, "question[") {
			err := f.updateQuestions(key, val)
			if err != nil {
				return f.getError(err.Error())
			}
		}

		if key == "title" {
			f.Title = val
		}

		if key == "description" {
			f.Description = val
		}
	}

	answers, err := f.gatherAnswers(postForm)
	if err != nil {
		return f.getError(err.Error())
	}

	f.updateAnswers(answers)

	return nil
}

func (f *editForm) updateQuestions(key, val string) error {
	var q questionEdit

	questionID, field, err := f.parseQuestionField(key)
	if err != nil {
		return err
	}

	if _, exists := f.Questions[questionID]; !exists {
		f.Questions[questionID] = q
	}

	q = f.Questions[questionID]

	switch field {
	case "content":
		q.Content = val
	case "correct":
		correct, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		q.Correct = correct
	case "points":
		points, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		q.Points = points
	}

	f.Questions[questionID] = q
	return nil
}

func (f *editForm) gatherAnswers(postForm url.Values) (answerEditMap, error) {
	answers := make(answerEditMap)

	for key, vals := range postForm {
		if !strings.HasPrefix(key, "answer[") {
			continue
		}

		answerID, field, err := f.parseAnswerField(key)
		if err != nil {
			return answers, err
		}

		if _, exists := answers[answerID]; !exists {
			answers[answerID] = answerEdit{}
		}

		a := answers[answerID]
		switch field {
		case "content":
			a.Content = vals[0]
		case "question":
			questionID, err := strconv.Atoi(vals[0])
			if err != nil {
				return answers, nil
			}
			a.ID = questionID
		}

		answers[answerID] = a
	}

	return answers, nil
}

func (f *editForm) updateAnswers(answers answerEditMap) {
	questionsToAnswers := make(answersEditMap)

	for answerID, answer := range answers {
		questionID := answer.ID

		if _, exists := questionsToAnswers[questionID]; !exists {
			questionsToAnswers[questionID] = []answerEdit{}
		}

		questionsToAnswers[questionID] = append(questionsToAnswers[questionID], answerEdit{
			ID:      answerID,
			Content: answer.Content,
		})
	}

	for questionID := range f.Questions {
		q := f.Questions[questionID]
		q.Answers = questionsToAnswers[questionID]
		f.Questions[questionID] = q
	}
}

func (f *editForm) parseQuestionField(fieldName string) (int, string, error) {
	id, field, err := parseQuestionField(fieldName, "question[id][content/question]")
	if err != nil {
		return id, field, f.getError(err.Error())
	}

	return id, field, err
}

func (f *editForm) parseAnswerField(fieldName string) (int, string, error) {
	re := regexp.MustCompile(`^answer\[(\d+)\]\[([^]]+)\]$`)
	matches := re.FindStringSubmatch(fieldName)
	if len(matches) != 3 {
		return 0, "", f.getError("answer name has insufficient data -- must match structure: answer[id][content/question]")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", f.getError(err.Error())
	}

	return id, matches[2], nil
}

func (f *editForm) getError(message string) error {
	return fmt.Errorf("edit form: malformed field (%s)", message)
}

func (f *editForm) serializeQuestionContent() map[string]string {
	s := make(map[string]string)

	for id, question := range f.Questions {
		s[fmt.Sprintf("question[%d][content]", id)] = question.Content
	}
	return s
}

func (f *editForm) serializeAnswerContent() map[string]string {
	s := make(map[string]string)

	for _, question := range f.Questions {
		for _, answer := range question.Answers {
			s[fmt.Sprintf("answer[%d][content]", answer.ID)] = answer.Content
		}
	}

	return s
}

func (f *editForm) serializeQuestionPoints() map[string]int {
	s := make(map[string]int)

	for id, question := range f.Questions {
		s[fmt.Sprintf("question[%d][points]", id)] = question.Points
	}

	return s
}

type questionAnswerMap map[int]int

type quizForm struct {
	quizID  int
	answers models.QuestionAnswer
}

func newQuizForm(postForm url.Values) (quizForm, error) {
	quizForm := quizForm{}
	err := quizForm.parseRequest(postForm)
	return quizForm, err
}

func (f *quizForm) parseRequest(postForm url.Values) error {
	f.answers = make(models.QuestionAnswer)

	for key, vals := range postForm {
		if len(vals) == 0 || key == "csrf_token" {
			continue
		}

		if strings.HasPrefix(key, "question[") {
			answerID, err := strconv.Atoi(vals[0])
			if err != nil {
				return err
			}

			questionID, err := f.parseQuestionField(key)
			if err != nil {
				return f.getError(err.Error())
			}

			f.answers[questionID] = answerID
		}

		if key == "quiz" {
			quizID, err := strconv.Atoi(vals[0])
			if err != nil {
				return f.getError("quiz ID must be an integer")
			}

			f.quizID = quizID
		}
	}

	return nil
}

func (f *quizForm) parseQuestionField(fieldName string) (int, error) {
	structure := "question[id][answer]"
	id, field, err := parseQuestionField(fieldName, structure)

	if err == nil && field != "answer" {
		err = fmt.Errorf("question has insufficient data -- must match structure: %s", structure)
	}

	return id, err
}

func parseQuestionField(fieldName, expectedStructure string) (int, string, error) {
	re := regexp.MustCompile(`^question\[(\d+)\]\[([^]]+)\]$`)
	matches := re.FindStringSubmatch(fieldName)
	if len(matches) != 3 {
		return 0, "", fmt.Errorf("question has insufficient data -- must match structure: %s", expectedStructure)
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", err
	}

	return id, matches[2], nil
}

func (f *quizForm) getError(message string) error {
	return fmt.Errorf("quiz form: malformed field (%s)", message)
}
