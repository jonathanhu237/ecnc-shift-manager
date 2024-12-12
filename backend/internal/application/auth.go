package application

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := app.readJSON(r, &payload); err != nil {
		app.errorResponse(w, r, app.badRequest(err.Error()))
		return
	}
	if err := app.validate.Struct(payload); err != nil {
		app.validateError(w, r, err)
		return
	}

	// get the user
	user, err := app.models.Users.SelectUserByUsername(payload.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.errorResponse(w, r, errInvalidLogin)
		default:
			app.internalSeverError(w, r, err)
		}
		return
	}

	// check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.errorResponse(w, r, errInvalidLogin)
		default:
			app.internalSeverError(w, r, err)
		}
		return
	}

	// create jwt
	claims := CustomClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(app.config.JWT.Secret))
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// generate refresh token
	refresh_token_string, err := app.generateRefreshToken()
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	refreshTokenExpiresAt := time.Now().Add(14 * 24 * time.Hour) // expires in 14 days

	refresh_token := &models.RefreshToken{
		UserID:           user.ID,
		RefreshTokenHash: app.hashRefreshToken(refresh_token_string),
		ExpiresAt:        refreshTokenExpiresAt,
	}

	// insert the refresh token to database
	if err := app.models.RefreshToken.InsertRefreshToken(refresh_token); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// set the refresh token in the http-only cookie
	cookie := &http.Cookie{
		Name:     "__ecnc_shift_manager_refresh_token",
		Value:    refresh_token_string,
		Path:     "/",
		Expires:  refreshTokenExpiresAt,
		HttpOnly: true,
	}

	if app.config.Environment == "production" {
		cookie.SameSite = http.SameSiteStrictMode
	}

	http.SetCookie(w, cookie)

	// response
	app.successResponse(w, r, "login successfully", map[string]any{
		"access_token": ss,
		"user":         user,
	})
}

func (app *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// get the requester ID
	requester, ok := r.Context().Value(requesterCtxKey).(*requester)
	if !ok {
		panic("requester not found in context")
	}

	// revoke user's refresh token
	if err := app.models.RefreshToken.RevokeUserTokens(requester.id); err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// remove the refresh token from the http-only cookie
	cookie := &http.Cookie{
		Name:    "__ecnc_shift_manager_refresh_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, cookie)

	// response
	app.successResponse(w, r, "logout successfully", nil)
}
