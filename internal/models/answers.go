package models

import (
	"database/sql"
	"time"
)

type Answer struct {
	ID      int
	Content string
	Correct bool
	Created time.Time
	Update  time.Time
	Deleted *time.Time
}

type AnswerPublic struct {
	Content string
	Correct bool
}

type AnswerModel struct {
	DB *sql.DB
}
