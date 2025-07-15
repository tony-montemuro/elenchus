package main

import "github.com/tony-montemuro/elenchus/internal/models"

type ProfilePageData struct {
	Published   []models.QuizMetadata
	Unpublished []models.QuizMetadata
}

type QuizzesPageData struct {
	Quizzes []models.QuizMetadata
}

type QuizPageData struct {
	Quiz models.QuizPublic
}

type AttemptPageData struct {
	Attempt models.AttemptPublic
}
