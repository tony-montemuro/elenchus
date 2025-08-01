package main

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"runtime/debug"
	"slices"
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

func (app *application) getProfileID(r *http.Request) (*int, error) {
	id := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)

	if id == 0 {
		return nil, errors.New("undefined profile id")
	}

	return &id, nil
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

func (app *application) redirectNotFound(w http.ResponseWriter, r *http.Request, logMessage string, err error) {
	app.redirectHome(w, r, "This page does not exist.", logMessage, err)
}

func (app *application) redirectHome(w http.ResponseWriter, r *http.Request, message, logMessage string, err error) {
	app.addErrorFlashToSession(message, r)
	app.logger.Warn(logMessage, slog.String("error", err.Error()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func pluralize(s string, n int) string {
	if n == 1 {
		return s
	}

	return fmt.Sprintf("%ss", s)
}

func sortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	slices.Sort(keys)
	return keys
}

func percentage(a, b int) string {
	x := float64(a) / float64(b)
	x = math.Round(x * 100)
	return fmt.Sprintf("%.2f%%", x)
}
