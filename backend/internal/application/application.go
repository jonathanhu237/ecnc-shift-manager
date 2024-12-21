package application

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	config    *config.Config
	logger    *slog.Logger
	validate  *validator.Validate
	server    *http.Server
	models    *models.Models
	emailChan *amqp.Channel
}

func New() *Application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &Application{
		logger:   logger,
		validate: validate,
	}
}

func (app *Application) Run() {
	// Read config
	cfg, err := config.ReadConfig(app.logger)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
	app.config = cfg

	// Establish database connection
	db, err := models.OpenDB(cfg)
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

	// establish rabbitmq mail producer
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

	// context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the mail sender
	if err := app.StartMailSender(ctx, ch); err != nil {
		app.logger.Error("failed to start the mail sender", slog.String("error", err.Error()))
		os.Exit(1)
	}
	app.logger.Info("email client established")

	// Start the server
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

func (app *Application) selfCheck() error {
	if err := app.checkBlackcoreExists(); err != nil {
		return err
	}

	return nil
}

func (app *Application) checkBlackcoreExists() error {
	exists, err := app.models.Users.CheckBlackcoreExists()
	if err != nil {
		return err
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(app.config.InitialAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if !exists {
		user := &models.User{
			Username:     app.config.InitialAdmin.Username,
			PasswordHash: string(password_hash), // ecnc_blackcore
			Email:        app.config.InitialAdmin.Email,
			FullName:     app.config.InitialAdmin.FullName,
			Role:         "黑心",
		}

		if err := app.models.Users.InsertUser(user); err != nil {
			return err
		}

		app.logger.Warn("blackcore does not exist, create a new one")
	}

	return nil
}
