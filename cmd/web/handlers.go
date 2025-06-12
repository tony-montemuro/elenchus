package main

import (
	"net/http"
)

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
