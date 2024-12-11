package application

import (
	"log/slog"
	"net/http"
	"time"
)

func (app *Application) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)

		app.logger.Info(
			"request processed",
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
			slog.Duration("duration", duration),
		)
	})
}
