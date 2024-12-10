package application

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

type Application struct {
	config *config.Config
	logger *slog.Logger
	server *http.Server
	models *models.Models
}

func New() *Application {
	cfg := config.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Application{
		config: cfg,
		logger: logger,
	}
}

func (app *Application) Run() {
	db, err := openDB(app.config)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app.logger.Info("database connection pool established")
	app.models = models.New(db)

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
