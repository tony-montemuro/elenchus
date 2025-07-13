package models

import "database/sql"

type QuestionAnswer map[int]int

type AttemptPublic struct {
	Quiz           QuizPublic
	PointsEarned   int
	QuestionAnswer QuestionAnswer
}

type AttemptModel struct {
	DB *sql.DB
}
