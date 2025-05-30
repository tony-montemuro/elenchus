package main

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

// func (c chain) thenFunc(h http.HandlerFunc) http.Handler {
// 	return c.then(h)
// }

func (c chain) then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}

	return h
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /ping", ping)

	globalChain := chain{app.recoverPanic, app.logRequest, commonHeaders}
	return globalChain.then(mux)
}
