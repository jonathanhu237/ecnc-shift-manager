package application

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment        string
	Port               int
	PostgresPassword   string
	JWTSecret          string
	MailClientSMTPHost string
	MailClientAddress  string
	MailClientPassword string
	RabbitMQPassword   string
}

func (app *Application) readConfig() (*Config, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.Environment = app.readOptionalStringEnv("ENVIRONMENT", "development")
	cfg.Port = app.readOptionalIntEnv("SERVER_PORT", 8080)

	postgresPassword, err := app.readRequiredStringEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.PostgresPassword = postgresPassword

	jwtSecret, err := app.readRequiredStringEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	cfg.JWTSecret = jwtSecret

	// mail client
	mailClientSMTPHost, err := app.readRequiredStringEnv("MAIL_CLIENT_SMTP_HOST")
	if err != nil {
		return nil, err
	}
	cfg.MailClientSMTPHost = mailClientSMTPHost

	mailClientAddress, err := app.readRequiredStringEnv("MAIL_CLIENT_ADDRESS")
	if err != nil {
		return nil, err
	}
	cfg.MailClientAddress = mailClientAddress

	mailClientPassword, err := app.readRequiredStringEnv("MAIL_CLIENT_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.MailClientPassword = mailClientPassword

	// rabbitmq
	rabbitMQPassword, err := app.readRequiredStringEnv("RABBITMQ_PASSWORD")
	if err != nil {
		return nil, err
	}
	cfg.RabbitMQPassword = rabbitMQPassword

	return cfg, nil
}

func (app *Application) readOptionalStringEnv(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return val
}

func (app *Application) readOptionalIntEnv(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		app.logger.Warn(
			"failed to parse environment variable, use default value",
			slog.String("key", key),
			slog.String("value", val),
			slog.Int("default_value", defaultValue),
		)
		return defaultValue
	}

	return intVal
}

func (app *Application) readRequiredStringEnv(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment %s is not set", key)
	}

	return val, nil
}
