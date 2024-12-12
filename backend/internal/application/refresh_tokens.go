package application

import (
	"errors"
	"net/http"
	"time"

	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) createRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	// get the user
	user, err := app.models.Users.SelectUserByUsername(payload.Username)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.unauthorized(w, r, errors.New("invalid username or password"))
		default:
			app.internalSeverError(w, r, err)
		}
		return
	}

	// check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.unauthorized(w, r, errors.New("invalid username or password"))
		default:
			app.internalSeverError(w, r, err)
		}
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
	if err := app.writeJSON(w, http.StatusCreated, nil); err != nil {
		app.internalSeverError(w, r, err)
	}
}
