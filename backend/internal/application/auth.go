package application

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := app.readJSON(r, &payload); err != nil {
		app.badRequest(w, r, err)
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
	expiresAt := time.Now().Add(24 * time.Hour) // expires in one day

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(app.config.JWTSecret))
	if err != nil {
		app.internalSeverError(w, r, err)
		return
	}

	// set the jwt in the http-only cookie
	cookie := &http.Cookie{
		Name:     "__ecnc_shift_manager_token",
		Value:    ss,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	}

	if app.config.Environment == "production" {
		cookie.SameSite = http.SameSiteStrictMode
	}

	http.SetCookie(w, cookie)

	// response
	app.successResponse(w, r, "login successfully", nil)
}

func (app *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// remove the token from the http-only cookie
	cookie := &http.Cookie{
		Name:    "__ecnc_shift_manager_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, cookie)

	// response
	app.successResponse(w, r, "logout successfully", nil)
}
