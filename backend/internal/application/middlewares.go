package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
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

func (app *Application) getRequesterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from cookie
		cookie, err := r.Cookie("__ecnc_shift_manager_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				app.errorResponse(w, r, errUnauthorized)

			default:
				app.internalSeverError(w, r, err)
			}
			return
		}

		// parse the token
		claims := &jwt.RegisteredClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(app.config.JWTSecret), nil
		})
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				app.errorResponse(w, r, errTokenIsExpired)
			default:
				app.errorResponse(w, r, errInvalidToken)
			}
			return
		}

		// get the requester details
		requesterUsername := claims.Subject
		requester, err := app.models.Users.SelectUserByUsername(requesterUsername)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.errorResponse(w, r, errInvalidToken)
			default:
				app.internalSeverError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), requesterCtxKey, requester)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	blackcoreLevel = 3
)

func (app *Application) authGuardMiddleware(levelRequired int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requester, ok := r.Context().Value(requesterCtxKey).(*requester)
			if !ok {
				panic("requester not found in context")
			}

			if requester.level < levelRequired {
				app.errorResponse(w, r, errForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) myInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requester, ok := r.Context().Value(requesterCtxKey).(*requester)
		if !ok {
			panic("myInfoMiddleware should be used after getUserInfoMiddleware")
		}

		details, err := app.models.Users.SelectUserByUsername(requester.username)
		if err != nil {
			app.internalSeverError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), requesterDetailsKey, details)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
