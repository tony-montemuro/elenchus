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

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl")
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "login.tmpl")
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "signup.tmpl")
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating a new user..."))
}
