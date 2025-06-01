package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/tony-montemuro/elenchus/internal/assert"
)

type LogOutput struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"msg"`
	Ip      string `json:"ip"`
	Proto   string `json:"proto"`
	Method  string `json:"method"`
	Uri     string `json:"uri"`
}

func TestCommonHeaders(t *testing.T) {
	rs := executeMiddleware(t, commonHeaders, http.HandlerFunc(ping))

	expectedValue := "default-src 'self'"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "same-origin"
	assert.Equal(t, rs.Header.Get("Cross-Origin-Opener-Policy"), expectedValue)

	expectedValue = "require-corp"
	assert.Equal(t, rs.Header.Get("Cross-Origin-Embedder-Policy"), expectedValue)

	expectedValue = "same-site"
	assert.Equal(t, rs.Header.Get("Cross-Origin-Resource-Policy"), expectedValue)

	expectedValue = "camera=(), geolocation=(), microphone=()"
	assert.Equal(t, rs.Header.Get("Permissions-Policy"), expectedValue)

	expectedValue = "max-age=6307200; includeSubDomains; preload"
	assert.Equal(t, rs.Header.Get("Strict-Transport-Security"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func TestRecoverPanic(t *testing.T) {
	middleware := newTestApplication(t, io.Discard).recoverPanic
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Panic handler")
	})

	rs := executeMiddleware(t, middleware, next)

	assert.Equal(t, rs.StatusCode, http.StatusInternalServerError)
	assert.Equal(t, rs.Header.Get("Connection"), "Closed")
}

func TestLogRequest(t *testing.T) {
	tests := []struct {
		name   string
		method string
		uri    string
		ip     string
	}{
		{
			name:   "GET request",
			method: http.MethodGet,
			uri:    "/",
			ip:     "127.0.0.1:8080",
		},
		{
			name:   "POST request",
			method: http.MethodPost,
			uri:    "/post",
			ip:     "248.247.196.20:4000",
		},
		{
			name:   "DELETE request",
			method: http.MethodDelete,
			uri:    "/delete",
			ip:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			middleware := newTestApplication(t, buf).logRequest
			handler := http.HandlerFunc(ping)

			rs := executeMiddlewareWithOptions(t, middleware, handler, tt.method, tt.uri, tt.ip)

			var output LogOutput
			err := json.Unmarshal(buf.Bytes(), &output)
			if err != nil {
				t.Fatal(err)
			}

			expectedValue := slog.LevelInfo.String()
			assert.Equal(t, output.Level, expectedValue)

			expectedValue = "receieved request"
			assert.Equal(t, output.Message, expectedValue)

			assert.Equal(t, output.Method, tt.method)
			assert.Equal(t, output.Uri, tt.uri)
			assert.Equal(t, output.Ip, tt.ip)

			expectedValue = "HTTP/1.1"
			assert.Equal(t, output.Proto, expectedValue)

			assert.Equal(t, rs.StatusCode, http.StatusOK)

		})
	}
}
