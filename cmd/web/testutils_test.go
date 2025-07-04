package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	modelMocks "github.com/tony-montemuro/elenchus/internal/models/mocks"
	serviceMocks "github.com/tony-montemuro/elenchus/internal/services/mocks"
)

func newTestApplication(t *testing.T, logWriter io.Writer) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		templateCache:  templateCache,
		logger:         slog.New(slog.NewJSONHandler(logWriter, nil)),
		profiles:       &modelMocks.ProfileModel{},
		quizzes:        &modelMocks.QuizModel{},
		quizzesService: &serviceMocks.QuizService{},
		sessionManager: sessionManager,
	}

	openAIClient := openai.NewClient(
		option.WithMaxRetries(0),
		option.WithMiddleware(app.LogOpenAIRequest),
	)

	app.openAIClient = openAIClient

	return app
}

func executeMiddleware(t *testing.T, middleware func(http.Handler) http.Handler, next http.Handler) *http.Response {
	return executeMiddlewareWithOptions(t, middleware, next, "GET", "/", "")
}

func executeMiddlewareWithOptions(t *testing.T, middleware func(http.Handler) http.Handler, next http.Handler, method, uri, ip string) *http.Response {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(method, uri, nil)
	if err != nil {
		t.Fatal(err)
	}

	r.RemoteAddr = ip

	middleware(next).ServeHTTP(rr, r)

	return rr.Result()
}

type testServer struct {
	*httptest.Server
}

func newTestServer(h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
