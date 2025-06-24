package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tony-montemuro/elenchus/internal/models"
	"github.com/tony-montemuro/elenchus/internal/validator"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	data.RangeRules = validator.RangeRules[validator.LoginForm]
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
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
		form.AddFieldError(err.Key, err.Error())
	}
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.RangeRules = rangeRules
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.profiles.Insert(form.FirstName, form.LastName, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use.")
			data := app.newTemplateData(r)
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
		form.AddFieldError(err.Key, err.Error())
	}
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")

	rangeRules := validator.RangeRules[formName]
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.RangeRules = rangeRules
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	profile, err := app.profiles.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect.")

			data := app.newTemplateData(r)
			data.Form = form
			data.RangeRules = rangeRules
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", profile.ID)
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Login successful! Welcome, %s!", profile.FirstName))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) quizList(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "quizzes.tmpl", data)
}
