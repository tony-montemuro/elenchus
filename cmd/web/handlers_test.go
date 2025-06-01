package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/tony-montemuro/elenchus/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t, io.Discard)

	ts := newTestServer(app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}
