package main

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

func (c chain) thenFunc(h http.HandlerFunc) http.Handler {
	return c.then(h)
}

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

	dynamicChain := chain{app.sessionManager.LoadAndSave, noSurf, app.authenticate}

	mux.Handle("GET /{$}", dynamicChain.thenFunc(app.home))
	mux.Handle("GET /login", dynamicChain.thenFunc(app.login))
	mux.Handle("POST /login", dynamicChain.thenFunc(app.loginPost))
	mux.Handle("GET /signup", dynamicChain.thenFunc(app.signup))
	mux.Handle("POST /signup", dynamicChain.thenFunc(app.signupPost))
	mux.Handle("POST /logout", dynamicChain.thenFunc(app.logoutPost))
	mux.Handle("GET /quizzes", dynamicChain.thenFunc(app.quizList))
	mux.Handle("GET /ping", dynamicChain.thenFunc(ping))

	protectedChain := append(dynamicChain, app.requireAuthentication)

	mux.Handle("GET /create", protectedChain.thenFunc(app.create))
	mux.Handle("POST /create", protectedChain.thenFunc(app.createPost))

	globalChain := chain{app.recoverPanic, app.logRequest, commonHeaders}
	return globalChain.then(mux)
}
