package handlers

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handlers struct {
	config    *config.Config
	logger    *slog.Logger
	models    *models.Models
	emailChan *amqp.Channel
	validate  *validator.Validate
}

func New(config *config.Config, logger *slog.Logger, models *models.Models, emailChan *amqp.Channel) *Handlers {
	return &Handlers{
		config:    config,
		logger:    logger,
		models:    models,
		emailChan: emailChan,
		validate:  validator.New(validator.WithRequiredStructEnabled()),
	}
}
