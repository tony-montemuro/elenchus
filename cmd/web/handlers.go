package main

import (
	"errors"
	"net/http"

	"github.com/tony-montemuro/elenchus/internal/models"
	"github.com/tony-montemuro/elenchus/internal/validator"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data.Form = loginForm{}
	data.RangeRules = validator.RangeRules[validator.LoginForm]
	data.Flash = flash
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = signupForm{}
	data.RangeRules = validator.RangeRules[validator.SignUpForm]
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := signupForm{
		FirstName: r.PostForm.Get("first-name"),
		LastName:  r.PostForm.Get("last-name"),
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
	}

	formName := validator.SignUpForm
	rangeRules := validator.RangeRules[formName]
	errs := validator.GetRangeErrors(form, formName)
	for _, err := range errs {
		form.AddError(err.Key, err.Error())
	}
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		data.RangeRules = rangeRules
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.profiles.Insert(form.FirstName, form.LastName, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddError("email", "Email address already in use.")
			data := app.newTemplateData()
			data.Form = form
			data.RangeRules = rangeRules
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful! Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := loginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	formName := validator.LoginForm
	errs := validator.GetRangeErrors(form, formName)
	for _, err := range errs {
		form.AddError(err.Key, err.Error())
	}

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		data.RangeRules = validator.RangeRules[formName]
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	w.Write([]byte("Logging in a user..."))
}
