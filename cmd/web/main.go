package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/gob"
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
	"github.com/tony-montemuro/elenchus/internal/ui"
)

type application struct {
	logger          *slog.Logger
	templateCache   map[string]*template.Template
	profiles        models.ProfileModelInterface
	quizzes         models.QuizModelInterface
	questionTypes   models.QuestionTypeModelInterface
	attempts        models.AttemptModelInterface
	quizzesService  services.QuizServiceInterface
	attemptsService services.AttemptServiceInterface
	sessionManager  *scs.SessionManager
	openAIClient    openai.Client
}

func init() {
	gob.Register(models.AttemptPublic{})
	gob.Register(profileForm{})
	gob.Register(createForm{})
	gob.Register(&ui.Flash{})
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
	attemptModel := &models.AttemptModel{DB: db}
	multipleChoiceAttemptModel := &models.MultipleChoiceAttemptModel{DB: db}
	quizService := &services.QuizService{
		DB:                         db,
		QuizModel:                  quizModel,
		QuestionModel:              questionModel,
		AnswerModel:                answerModel,
		MultipleChoiceAttemptModel: multipleChoiceAttemptModel,
		AttemptModel:               attemptModel,
		QuestionTypeModel:          questionTypeModel,
	}

	app := &application{
		logger:         logger,
		templateCache:  templateCache,
		profiles:       &models.ProfileModel{DB: db},
		quizzes:        quizModel,
		questionTypes:  questionTypeModel,
		attempts:       attemptModel,
		quizzesService: quizService,
		attemptsService: &services.AttemptService{
			DB:                         db,
			AttemptModel:               attemptModel,
			MultipleChoiceAttemptModel: multipleChoiceAttemptModel,
			QuizService:                quizService,
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
