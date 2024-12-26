package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/handlers"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/utils"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/workers"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Application struct {
	config    *config.Config
	logger    *slog.Logger
	server    *http.Server
	handler   *handlers.Handlers
	models    *models.Models
	emailChan *amqp.Channel
}

func New() *Application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Application{
		logger: logger,
	}
}

func (app *Application) Run() {
	/****************************************************************
		read config
	****************************************************************/

	cfg, err := config.ReadConfig(app.logger)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.config = cfg

	/****************************************************************
		establish database connection
	****************************************************************/
	db, err := utils.OpenDB(cfg)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	app.logger.Info("database connection pool established")
	app.models = models.New(db)

	/****************************************************************
		establish mail sender
	****************************************************************/
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", app.config.RabbitMQ.User, app.config.RabbitMQ.Password, app.config.RabbitMQ.Host))
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	defer ch.Close()
	app.emailChan = ch

	_, err = ch.QueueDeclare("mail_queue", true, false, false, false, nil)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background()) // context for graceful shutdown
	defer cancel()

	mailSender := workers.NewMailSender(app.config, app.logger, ch)
	if err := mailSender.Run(ctx); err != nil {
		app.logger.Error("failed to start the mail sender", slog.String("error", err.Error()))
		os.Exit(1)
	}
	app.logger.Info("email client established")

	/****************************************************************
		perform mail sender
	****************************************************************/
	if err := app.healthCheck(); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.logger.Info("health check completed")

	/****************************************************************
		establish mail sender
	****************************************************************/
	app.handler = handlers.New(app.config, app.logger, app.models, ch)
	app.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.ServerPort),
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
