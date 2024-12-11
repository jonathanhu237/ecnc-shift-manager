package config

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

type Config struct {
	Environment string
	Server      Server
	Database    Database
	JWT         JWT
}

func New() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Environment, "environment", "development", "Application environment")
	flag.IntVar(&cfg.Server.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Database.DSN, "db-dsn", os.Getenv("ECNC_SHIFT_MANAGER_DB_DSN"), "PostgreSQL DSN")
	flag.StringVar(&cfg.JWT.Secret, "jwt-secret", os.Getenv("ECNC_SHIFT_MANAGER_JWT_SECRET"), "JWT secret key")
	flag.Parse()

	return cfg
}
