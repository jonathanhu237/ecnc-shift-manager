package application

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
)

type Application struct {
	config *config.Config
	logger *slog.Logger
	server *http.Server
}

func New() *Application {
	cfg := config.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Application{
		config: cfg,
		logger: logger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
			Handler:      routes(),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		},
	}
}

func (app *Application) Run() {
	app.logger.Info("starting server", "addr", app.server.Addr)
	if err := app.server.ListenAndServe(); err != nil {
		app.logger.Error(err.Error())
	}
}
