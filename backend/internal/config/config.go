package config

import (
	"flag"
	"time"
)

type Server struct {
	Port         int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Config struct {
	Server Server
}

func New() *Config {
	cfg := &Config{}

	flag.IntVar(&cfg.Server.Port, "port", 8080, "API server port")
	flag.DurationVar(&cfg.Server.IdleTimeout, "idle-timeout", time.Minute, "API server idle timeout")
	flag.DurationVar(&cfg.Server.ReadTimeout, "read-timeout", 5*time.Second, "API server read timeout")
	flag.DurationVar(&cfg.Server.WriteTimeout, "write-timeout", 10*time.Second, "API server write timeout")
	flag.Parse()

	return cfg
}
