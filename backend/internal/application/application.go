package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

type Application struct {
	config   *Config
	logger   *slog.Logger
	server   *http.Server
	models   *models.Models
	validate *validator.Validate
}

func New() *Application {
	cfg := readConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Application{
		config:   cfg,
		logger:   logger,
		validate: validate,
	}
}

func (app *Application) Run() {
	db, err := app.openDB()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app.logger.Info("database connection pool established")
	app.models = models.New(db)

	if err := app.selfCheck(); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.logger.Info("self check completed")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.StartTokenCleaner(ctx)
	app.logger.Info("token cleaner started")

	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Server.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	app.logger.Info("starting server", "addr", app.server.Addr)
	if err := app.server.ListenAndServe(); err != nil {
		app.logger.Error(err.Error())
	}
}
