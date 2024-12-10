package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
)

func New(cfg *config.Config, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
