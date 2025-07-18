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
	Attempts  []models.AttemptMetadata
	Quiz      models.QuizPublic
	ProfileID int
}

func (d *QuizPageData) setProfileID(profileID *int) {
	if profileID != nil {
		d.ProfileID = *profileID
	} else {
		d.ProfileID = 0
	}
}

type AttemptPageData struct {
	Attempt models.AttemptPublic
}
