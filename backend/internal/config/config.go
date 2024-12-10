package config

import (
	"flag"
)

type Server struct {
	Port int
}

type Config struct {
	Server Server
}

func New() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.Server.Port, "port", 8080, "API server port")
	flag.Parse()

	return cfg
}
