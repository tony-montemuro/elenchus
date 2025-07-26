package main

import (
	"errors"
	"net/http"
)

var ErrNoCreateForm = errors.New("session: no create form found")

func (app *application) redirectToCreate(w http.ResponseWriter, r *http.Request, form createForm) {
	app.sessionManager.Put(r.Context(), createFormKey, form)
	redirectPath := "/create?view=" + form.Type
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func (app *application) getCreateFormFromSession(r *http.Request) (createForm, error) {
	var form createForm
	if !app.sessionManager.Exists(r.Context(), createFormKey) {
		return form, ErrNoCreateForm
	}

	form, ok := app.sessionManager.Pop(r.Context(), createFormKey).(createForm)
	if !ok {
		return form, errors.New("session profile form not correctly typed! must be of type `createForm`")
	}

	return form, nil
}
