package main

import "github.com/tony-montemuro/elenchus/internal/models"

type ProfilePageData struct {
	Published   []models.QuizMetadata
	Unpublished []models.QuizMetadata
}
