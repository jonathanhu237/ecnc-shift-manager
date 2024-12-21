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
	}

	RabbitMQ struct {
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

	cfg.Environment = cfg.readStringEnv("ENVIRONMENT", "development")
	cfg.ServerPort = cfg.readIntEnv("SERVER_PORT", 8080)
	cfg.JWTSecret = cfg.readStringEnv("JWT_SECRET", "")

	// postgres
	cfg.Postgres.User = cfg.readStringEnv("POSTGRES_USER", "postgres")
	cfg.Postgres.Password = cfg.readStringEnv("POSTGRES_PASSWORD", "")
	cfg.Postgres.DB = cfg.readStringEnv("POSTGRES_DB", "ecnc_shift_manager_db")

	// rabbitmq
	cfg.RabbitMQ.User = cfg.readStringEnv("RABBITMQ_DEFAULT_USER", "rabbitmq")
	cfg.RabbitMQ.Password = cfg.readStringEnv("RABBITMQ_DEFAULT_PASS", "")

	// Email
	cfg.MailClient.SMTPHost = cfg.readStringEnv("MAIL_CLIENT_SMTP_HOST", "")
	cfg.MailClient.Sender = cfg.readStringEnv("MAIL_CLIENT_SENDER", "")
	cfg.MailClient.Password = cfg.readStringEnv("MAIL_CLIENT_PASSWORD", "")

	return cfg, nil
}

func (cfg *Config) readStringEnv(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		cfg.logger.Warn(
			"Environment variable not found, use default value instead.",
			slog.String("key", key),
			slog.String("default_value", defaultValue),
		)
		return defaultValue
	}

	return val
}

func (cfg *Config) readIntEnv(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		cfg.logger.Warn(
			"Environment variable not found, use default value instead.",
			slog.String("key", key),
			slog.Int("default_value", defaultValue),
		)
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		cfg.logger.Warn(
			"Environment variable is not a valid integer, use default value instead.",
			slog.String("key", key),
			slog.String("value", val),
			slog.Int("default_value", defaultValue),
		)
		return defaultValue
	}

	return intVal
}
