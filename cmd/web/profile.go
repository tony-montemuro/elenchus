package main

import (
	"errors"
	"net/http"
)

var ErrNoProfileForm = errors.New("session: no profile form found")

func (app *application) getProfileFormFromSession(r *http.Request) (profileForm, error) {
	var form profileForm
	if !app.sessionManager.Exists(r.Context(), profileFormKey) {
		return form, ErrNoProfileForm
	}

	form, ok := app.sessionManager.Pop(r.Context(), profileFormKey).(profileForm)
	if !ok {
		return form, errors.New("session profile form not correctly typed! must be of type `profileForm`")
	}

	return form, nil
}
