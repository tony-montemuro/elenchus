package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.tmpl")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
