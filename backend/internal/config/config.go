package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	logger *slog.Logger

	Environment        string
	Port               int
	PostgresPassword   string
	JWTSecret          string
	MailClientSMTPHost string
	MailClientAddress  string
	MailClientPassword string
	RabbitMQPassword   string
}

func ReadConfig(logger *slog.Logger) (*Config, error) {
	cfg := &Config{
		logger: logger,
	}

	cfg.Environment = cfg.readOptionalStringEnv("ENVIRONMENT", "development")
	cfg.Port = cfg.readOptionalIntEnv("SERVER_PORT", 8080)

	postgresPassword, err := cfg.readRequiredStringEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.PostgresPassword = postgresPassword

	jwtSecret, err := cfg.readRequiredStringEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	cfg.JWTSecret = jwtSecret

	// mail client
	mailClientSMTPHost, err := cfg.readRequiredStringEnv("MAIL_CLIENT_SMTP_HOST")
	if err != nil {
		return nil, err
	}
	cfg.MailClientSMTPHost = mailClientSMTPHost

	mailClientAddress, err := cfg.readRequiredStringEnv("MAIL_CLIENT_ADDRESS")
	if err != nil {
		return nil, err
	}
	cfg.MailClientAddress = mailClientAddress

	mailClientPassword, err := cfg.readRequiredStringEnv("MAIL_CLIENT_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.MailClientPassword = mailClientPassword

	// rabbitmq
	rabbitMQPassword, err := cfg.readRequiredStringEnv("RABBITMQ_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.RabbitMQPassword = rabbitMQPassword

	return cfg, nil
}

func (cfg *Config) readOptionalStringEnv(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return val
}

func (cfg *Config) readOptionalIntEnv(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		cfg.logger.Warn(
			"failed to parse environment variable, use default value",
			slog.String("key", key),
			slog.String("value", val),
			slog.Int("default_value", defaultValue),
		)
		return defaultValue
	}

	return intVal
}

func (cfg *Config) readRequiredStringEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment %s is not set", key)
	}

	return val, nil
}
