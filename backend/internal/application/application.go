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
	"github.com/wneessen/go-mail"
)

type Application struct {
	config     *Config
	logger     *slog.Logger
	validate   *validator.Validate
	server     *http.Server
	models     *models.Models
	mailClient *mail.Client
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
	// Establish database connection
	db, err := app.openDB()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app.logger.Info("database connection pool established")
	app.models = models.New(db)

	// Perform self check
	if err := app.selfCheck(); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.logger.Info("self check completed")

	// Start token cleaner
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app.StartTokenCleaner(ctx)
	app.logger.Info("token cleaner started")

	// Establish email client
	app.mailClient, err = mail.NewClient(
		"smtp.feishu.cn",
		mail.WithPort(465),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithSSL(),
		mail.WithUsername(app.config.Email.Address),
		mail.WithPassword(app.config.Email.Password),
	)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.mailClient.DialWithContext(ctx); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	app.logger.Info("email client established")

	// Start the server
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
