package main

import (
	"errors"
	"net/http"

	"github.com/tony-montemuro/elenchus/internal/models"
)

var ErrNoAttempt = errors.New("session: no attempt found")

func (app *application) getAttemptFromSession(r *http.Request) (models.AttemptPublic, error) {
	var attempt models.AttemptPublic
	if !app.sessionManager.Exists(r.Context(), attemptKey) {
		return attempt, ErrNoAttempt
	}

	attempt, ok := app.sessionManager.Pop(r.Context(), attemptKey).(models.AttemptPublic)
	if !ok {
		return attempt, errors.New("session attempt not correctly typed! must be of type `models.AttemptPublic`")
	}

	return attempt, nil
}
