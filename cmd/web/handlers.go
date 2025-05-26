package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("/ route")
	w.Write([]byte("Hello world!"))
}
