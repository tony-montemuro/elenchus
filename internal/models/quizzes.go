package models

import (
	"database/sql"
	"time"
)

type QuizModelInterface interface {
	Latest() ([]QuizMetadata, error)
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
	Profile       Profile
	Title         string
	Description   string
	QuestionCount int
	Published     *time.Time
}

type QuizModel struct {
	DB *sql.DB
}

func (m *QuizModel) Latest() ([]QuizMetadata, error) {
	return []QuizMetadata{}, nil
}
