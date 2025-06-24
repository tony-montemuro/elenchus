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

type QuestionModel struct {
	DB *sql.DB
}
