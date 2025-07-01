package models

import (
	"database/sql"
	"time"
)

type Question struct {
	ID      int
	Type    QuestionType
	Answers []Answer
	Content string
	Points  uint32
	Created time.Time
	Update  time.Time
	Deleted *time.Time
}

type QuestionPublic struct {
	ID      int
	Type    QuestionTypePublic
	Content string
	Points  int
	Answers []AnswerPublic
}

type QuestionModel struct {
	DB *sql.DB
}
