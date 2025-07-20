package main

import (
	"net/http"
	"slices"

	"github.com/tony-montemuro/elenchus/ui"
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

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamicChain := chain{app.sessionManager.LoadAndSave, noSurf, app.authenticate}

	mux.Handle("GET /{$}", dynamicChain.thenFunc(app.home))
	mux.Handle("GET /login", dynamicChain.thenFunc(app.login))
	mux.Handle("POST /login", dynamicChain.thenFunc(app.loginPost))
	mux.Handle("GET /signup", dynamicChain.thenFunc(app.signup))
	mux.Handle("POST /signup", dynamicChain.thenFunc(app.signupPost))
	mux.Handle("POST /logout", dynamicChain.thenFunc(app.logoutPost))
	mux.Handle("GET /quizzes", dynamicChain.thenFunc(app.quizList))
	mux.Handle("GET /quizzes/{quizID}", dynamicChain.thenFunc(app.quiz))
	mux.Handle("POST /quizzes/{quizID}", dynamicChain.thenFunc(app.quizPost))
	mux.Handle("GET /quizzes/{quizID}/result", dynamicChain.thenFunc(app.result))
	mux.Handle("GET /profile/{profileID}", dynamicChain.thenFunc(app.profile))
	mux.Handle("GET /ping", dynamicChain.thenFunc(ping))

	protectedChain := append(dynamicChain, app.requireAuthentication)

	mux.Handle("GET /create", protectedChain.thenFunc(app.create))
	mux.Handle("POST /create", protectedChain.thenFunc(app.createPost))
	mux.Handle("GET /profile", protectedChain.thenFunc(app.myProfile))
	mux.Handle("POST /profile", protectedChain.thenFunc(app.myProfilePost))
	mux.Handle("GET /quizzes/{quizID}/edit", protectedChain.thenFunc(app.edit))
	mux.Handle("POST /quizzes/{quizID}/edit", protectedChain.thenFunc(app.editPost))
	mux.Handle("GET /quizzes/{quizID}/attempt/{attemptID}", protectedChain.thenFunc(app.attempt))
	mux.Handle("POST /quizzes/{quizID}/unpublish", protectedChain.thenFunc(app.unpublish))

	globalChain := chain{app.recoverPanic, app.addRequestID, app.logRequest, commonHeaders}
	return globalChain.then(mux)
}
