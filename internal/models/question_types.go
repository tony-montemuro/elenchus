package models

import (
	"database/sql"
	"time"
)

type QuestionType struct {
	ID            int
	Name          string
	DefaultPoints uint32
	Created       time.Time
	Update        time.Time
}

type QuestionTypeModel struct {
	DB *sql.DB
}
