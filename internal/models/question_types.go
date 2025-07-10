package models

import (
	"database/sql"
	"time"
)

type QuestionTypeModelInterface interface {
	GetMultipleChoiceId() (int, error)
	GetMultipleChoice() (QuestionTypePublic, error)
}

type QuestionType struct {
	ID            int
	Name          string
	DefaultPoints uint32
	Created       time.Time
	Update        time.Time
}

type QuestionTypePublic struct {
	Name          string
	DefaultPoints uint32
}

type QuestionTypeModel struct {
	DB *sql.DB
}

var (
	MultipleChoice = "multiple_choice"
	FreeResponse   = "free_response"
)

func (m *QuestionTypeModel) GetMultipleChoiceId() (int, error) {
	var id int

	stmt := `SELECT id
	FROM question_type
	WHERE name = ?`

	err := m.DB.QueryRow(stmt, MultipleChoice).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *QuestionTypeModel) GetMultipleChoice() (QuestionTypePublic, error) {
	var qt QuestionTypePublic
	var err error

	stmt := `SELECT name, default_points 
	FROM question_type
	WHERE name = ?`

	err = m.DB.QueryRow(stmt, MultipleChoice).Scan(&qt.Name, &qt.DefaultPoints)
	return qt, err
}
