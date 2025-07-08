package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

type questionEdit struct {
	Content string
	Correct int
}

type answerEdit struct {
	Content string
}

type questionEditMap map[int]questionEdit
type answerEditMap map[int]answerEdit

type editForm struct {
	Title       string
	Description string
	Questions   questionEditMap
	Answers     answerEditMap
	validator.Validator
}

func (f *editForm) parseRequest(postForm url.Values) error {
	f.Questions = make(questionEditMap)
	f.Answers = make(answerEditMap)

	for key, vals := range postForm {
		if len(vals) == 0 || key == "csrf_token" {
			continue
		}

		if strings.HasPrefix(key, "question[") {
			questionID, field, err := f.parseQuestionField(key)
			if err != nil {
				return f.getError(err.Error())
			}

			if _, exists := f.Questions[questionID]; !exists {
				f.Questions[questionID] = questionEdit{}
			}

			q := f.Questions[questionID]
			switch field {
			case "content":
				q.Content = vals[0]
			case "correct":
				correct, err := strconv.Atoi(vals[0])
				if err != nil {
					return f.getError(err.Error())
				}
				q.Correct = correct
			}
			f.Questions[questionID] = q
		}

		if strings.HasPrefix(key, "answer[") {
			answerID, field, err := f.parseAnswerField(key)
			if err != nil {
				return f.getError(err.Error())
			}

			if field != "content" {
				return f.getError("answer has bad structure -- second key must be 'content'")
			}

			if _, exists := f.Answers[answerID]; !exists {
				f.Answers[answerID] = answerEdit{}
			}

			answer := f.Answers[answerID]
			answer.Content = vals[0]
			f.Answers[answerID] = answer
		}

		if key == "title" {
			f.Title = vals[0]
		}

		if key == "description" {
			f.Description = vals[0]
		}
	}

	return nil
}

func (f *editForm) parseQuestionField(fieldName string) (int, string, error) {
	re := regexp.MustCompile(`^question\[(\d+)\]\[([^]]+)\]$`)
	matches := re.FindStringSubmatch(fieldName)
	if len(matches) != 3 {
		return 0, "", f.getError("question has insufficient data -- must match structure: question[id][content/correct]")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", f.getError(err.Error())
	}

	return id, matches[2], nil
}

func (f *editForm) parseAnswerField(fieldName string) (int, string, error) {
	re := regexp.MustCompile(`^answer\[(\d+)\]\[([^]]+)\]$`)
	matches := re.FindStringSubmatch(fieldName)
	if len(matches) != 3 {
		return 0, "", f.getError("answer name has insufficient data -- must match structure: answer[id][content]")
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, "", f.getError(err.Error())
	}

	return id, matches[2], nil
}

func (f *editForm) getError(message string) error {
	return fmt.Errorf("edit form: malformed answer field (%s)", message)
}
