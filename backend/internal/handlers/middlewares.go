package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
)

func (h *Handlers) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)

		h.logger.Info(
			"request processed",
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
			slog.Duration("duration", duration),
		)
	})
}

func (h *Handlers) GetRequesterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from cookie
		cookie, err := r.Cookie("__ecnc_shift_manager_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				h.errorResponse(w, r, errors.New("用户未登录"))
				return
			default:
				h.internalServerError(w, r, err)
				return
			}
		}

		// parse the token
		claims := &jwt.RegisteredClaims{}
		_, err = jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(h.config.JWTSecret), nil
		})
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				h.errorResponse(w, r, errors.New("无效的访问令牌"))
				return
			default:
				h.internalServerError(w, r, err)
				return
			}
		}

		// get the requester details
		requesterUsername := claims.Subject
		requester, err := h.models.SelectUserByUsername(requesterUsername)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				h.errorResponse(w, r, errors.New("无效的访问令牌"))
				return
			default:
				h.internalServerError(w, r, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), requesterCtxKey, requester)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handlers) AuthGuardMiddleware(levelRequired int32) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requester, ok := r.Context().Value(requesterCtxKey).(*models.User)
			if !ok {
				panic("AuthGuardMiddleware must used after GetRequesterMiddleware")
			}

			if requester.Level < levelRequired {
				h.errorResponse(w, r, errors.New("权限不足"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handlers) GetUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDParam := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(userIDParam, 10, 64)
		if err != nil {
			h.errorResponse(w, r, errors.New("无效的用户ID"))
			return
		}

		user, err := h.models.SelectUserByID(userID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				h.errorResponse(w, r, errors.New("用户不存在"))
			default:
				h.internalServerError(w, r, err)
				return
			}
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
