package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
	"github.com/openai/openai-go/option"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		w.Header().Set("Permissions-Policy", "camera=(), geolocation=(), microphone=()")
		w.Header().Set("Strict-Transport-Security", "max-age=6307200; includeSubDomains; preload")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequestInfo(r *http.Request, requestID string, message string, level slog.Level) {
	app.logger.LogAttrs(r.Context(), level, message,
		slog.String("request_id", requestID),
		slog.String("ip", r.RemoteAddr),
		slog.String("proto", r.Proto),
		slog.String("method", r.Method),
		slog.String("uri", r.URL.RequestURI()),
	)
}

func (app *application) logResponseInfo(res *http.Response, err error, requestID string, elapsed time.Duration, message string, level slog.Level) {
	attrs := []slog.Attr{
		slog.String("request_id", requestID),
		slog.Any("error", err),
		slog.Int64("time_ms", elapsed.Milliseconds()),
	}

	if res != nil {
		attrs = append(attrs,
			slog.Int("status_code", res.StatusCode),
			slog.String("status", res.Status),
			slog.String("proto", res.Proto),
		)
	}

	app.logger.LogAttrs(context.Background(), level, message, attrs...)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logRequestInfo(r, getRequestID(r), "receieved request", slog.LevelInfo)

		next.ServeHTTP(w, r)
	})
}

func (app *application) addRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Closed")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.logRequestInfo(r, getRequestID(r), "unauthenticated user attempted to access a protected route", slog.LevelWarn)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.profiles.Exists(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) LogOpenAIRequest(r *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	start := time.Now()
	requestID := getRequestID(r)
	app.logRequestInfo(r, requestID, "Begin OpenAI Request", slog.LevelInfo)

	res, err = next(r)

	app.logResponseInfo(res, err, requestID, time.Since(start), "Stop OpenAI Response", slog.LevelInfo)
	return res, err
}
