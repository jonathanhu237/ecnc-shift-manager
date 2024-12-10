package application

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/server"
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
		server: server.New(cfg, logger),
	}
}

func (app *Application) Run() {
	app.logger.Info("starting server", "addr", app.server.Addr)
	if err := app.server.ListenAndServe(); err != nil {
		app.logger.Error(err.Error())
	}
}
