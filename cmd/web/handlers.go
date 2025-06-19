package main

import (
	"net/http"

	"github.com/tony-montemuro/elenchus/internal/validator"
)

type signupForm struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	validator.Validator
}

func (f signupForm) GetStringVals() map[string]string {
	vals := make(map[string]string)

	vals["firstName"] = f.FirstName
	vals["lastName"] = f.LastName
	vals["email"] = f.Email
	vals["password"] = f.Password

	return vals
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = signupForm{}
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

	err = validator.InputsInRange(form, validator.SignUpForm)
	if err != nil {
		app.logger.Warn(err.Error())
	}

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		data.RangeRules = validator.RangeRules[validator.SignUpForm]
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	w.Write([]byte("Creating a new user..."))
}
