package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/tony-montemuro/elenchus/internal/config"
)

type application struct {
	logger *slog.Logger
}

func main() {
	config := config.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.MinLogLevel,
	}))
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *config.Addr), slog.String("minLoggingLevel", config.MinLogLevel.String()))
	srv := &http.Server{
		Addr:     *config.Addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
