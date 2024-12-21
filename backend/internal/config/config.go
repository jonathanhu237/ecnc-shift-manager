package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	logger *slog.Logger

	Environment string
	ServerPort  int
	JWTSecret   string

	Postgres struct {
		User     string
		Password string
		DB       string
		Host     string
	}

	RabbitMQ struct {
		Host     string
		User     string
		Password string
	}

	MailClient struct {
		SMTPHost string
		Sender   string
		Password string
	}
}

func ReadConfig(logger *slog.Logger) (*Config, error) {
	cfg := &Config{
		logger: logger,
	}

	cfg.Environment = cfg.readStringEnv("ENVIRONMENT")
	cfg.ServerPort = cfg.readIntEnv("API_SERVER_PORT")
	cfg.JWTSecret = cfg.readStringEnv("JWT_SECRET")

	// postgres
	cfg.Postgres.User = cfg.readStringEnv("POSTGRES_USER")
	cfg.Postgres.Password = cfg.readStringEnv("POSTGRES_PASSWORD")
	cfg.Postgres.DB = cfg.readStringEnv("POSTGRES_DB")
	cfg.Postgres.Host = cfg.readStringEnv("POSTGRES_HOST")

	// rabbitmq
	cfg.RabbitMQ.Host = cfg.readStringEnv("RABBITMQ_HOST")
	cfg.RabbitMQ.User = cfg.readStringEnv("RABBITMQ_DEFAULT_USER")
	cfg.RabbitMQ.Password = cfg.readStringEnv("RABBITMQ_DEFAULT_PASS")

	// Email
	cfg.MailClient.SMTPHost = cfg.readStringEnv("MAIL_CLIENT_SMTP_HOST")
	cfg.MailClient.Sender = cfg.readStringEnv("MAIL_CLIENT_SENDER")
	cfg.MailClient.Password = cfg.readStringEnv("MAIL_CLIENT_PASSWORD")

	return cfg, nil
}

func (cfg *Config) readStringEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		cfg.logger.Warn("Environment variable is empty.", slog.String("key", key))
	}
	return val
}

func (cfg *Config) readIntEnv(key string) int {
	val := os.Getenv(key)
	if val == "" {
		cfg.logger.Warn("Environment variable is empty, use zero instead.", slog.String("key", key))
		return 0
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		cfg.logger.Warn(
			"Environment variable is not a valid integer, use zero instead.",
			slog.String("key", key),
			slog.String("value", val),
		)
		return 0
	}

	return intVal
}
