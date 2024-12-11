package application

import (
	"context"
	"log/slog"
	"time"
)

func (app *Application) StartTokenCleaner(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		if err := app.models.RefreshToken.DeleteExpiredTokens(); err != nil {
			app.logger.Error("failed to delete expired refresh tokens", slog.String("error", err.Error()))
		}

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := app.models.RefreshToken.DeleteExpiredTokens(); err != nil {
					app.logger.Error("failed to delete expired refresh tokens", slog.String("error", err.Error()))
				}
			}
		}
	}()
}
