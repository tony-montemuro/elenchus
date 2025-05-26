package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	app.logger.Info("home route")
	w.Write([]byte("Hello world!"))
}
