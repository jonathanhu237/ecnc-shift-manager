package handlers

import (
	"log/slog"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handlers struct {
	config    *config.Config
	logger    *slog.Logger
	models    *models.Models
	emailChan *amqp.Channel
}

func New(config *config.Config, logger *slog.Logger, models *models.Models, emailChan *amqp.Channel) *Handlers {
	return &Handlers{
		config:    config,
		logger:    logger,
		models:    models,
		emailChan: emailChan,
	}
}
