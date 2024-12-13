package application

import (
	"flag"
	"os"
)

type Server struct {
	Port int
}

type Database struct {
	DSN string
}

type JWT struct {
	Secret string
}

type Email struct {
	Address  string
	Password string
}

type Config struct {
	Environment string
	Server      Server
	Database    Database
	JWT         JWT
	Email       Email
}

func readConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Environment, "environment", "development", "Application environment")
	flag.IntVar(&cfg.Server.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Database.DSN, "db-dsn", os.Getenv("ECNC_SHIFT_MANAGER_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", os.Getenv("ECNC_SHIFT_MANAGER_JWT_SECRET"), "JWT secret key")
	flag.StringVar(&cfg.Email.Address, "email-sender", os.Getenv("ECNC_SHIFT_MANAGER_EMAIL_ADDRESS"), "Email sender address")
	flag.StringVar(&cfg.Email.Password, "email-sender-password", os.Getenv("ECNC_SHIFT_MANAGER_EMAIL_PASSWORD"), "Email sender password")
	flag.Parse()

	return cfg
}
