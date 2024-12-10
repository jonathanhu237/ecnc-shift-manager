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

type Config struct {
	Server   Server
	Database Database
}

func New() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.Server.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Database.DSN, "db-dsn", os.Getenv("ECNC_SHIFT_MANAGER_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	return cfg
}
