package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (app *Application) getUserInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.errorResponse(w, r, errAuthHeaderNotSet)
			return
		}

		// check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			app.errorResponse(w, r, errInvalidAuthHeader)
			return
		}

		// extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// parse the token
		claims := &CustomClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
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

		requesterID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorResponse(w, r, errInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), requesterCtxKey, &requester{
			id:    requesterID,
			role:  claims.Role,
			level: claims.Level,
		})

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
