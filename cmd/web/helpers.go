package main

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("stack trace", trace))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func generateRequestID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(requestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
