package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
)

func New(cfg *config.Config, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes(),
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
