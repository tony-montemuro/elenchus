package main

import (
	"crypto/tls"
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tony-montemuro/elenchus/internal/config"
	"github.com/tony-montemuro/elenchus/internal/models"
	"github.com/tony-montemuro/elenchus/internal/services"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	profiles       models.ProfileModelInterface
	quizzes        models.QuizModelInterface
	questionTypes  models.QuestionTypeModelInterface
	quizzesService services.QuizServiceInterface
	sessionManager *scs.SessionManager
	openAIClient   openai.Client
}

func main() {
	config := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.MinLogLevel,
	}))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDB(*config.Dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	defer db.Close()

	quizModel := &models.QuizModel{DB: db}
	questionModel := &models.QuestionModel{DB: db}
	answerModel := &models.AnswerModel{DB: db}
	questionTypeModel := &models.QuestionTypeModel{DB: db}

	app := &application{
		logger:        logger,
		templateCache: templateCache,
		profiles:      &models.ProfileModel{DB: db},
		quizzes:       quizModel,
		questionTypes: questionTypeModel,
		quizzesService: &services.QuizService{
			QuizModel:         quizModel,
			QuestionModel:     questionModel,
			AnswerModel:       answerModel,
			QuestionTypeModel: questionTypeModel,
		},
		sessionManager: sessionManager,
	}

	openAIClient := openai.NewClient(
		option.WithMaxRetries(0),
		option.WithMiddleware(app.LogOpenAIRequest),
	)

	app.openAIClient = openAIClient

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	logger.Info("starting server", slog.String("addr", *config.Addr), slog.String("minLoggingLevel", config.MinLogLevel.String()))
	srv := &http.Server{
		Addr:      *config.Addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
	}

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
