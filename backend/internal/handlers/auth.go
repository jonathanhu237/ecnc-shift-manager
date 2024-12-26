package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := h.readJSON(r, &payload); err != nil {
		h.errorResponse(w, r, err)
		return
	}

	// check the payload
	switch {
	case payload.Username == "":
		h.errorResponse(w, r, errors.New("用户名为空"))
		return
	case payload.Password == "":
		h.errorResponse(w, r, errors.New("密码为空"))
		return
	}

	// get the user
	user, err := h.models.SelectUserByUsername(payload.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.errorResponse(w, r, errors.New("用户名不存在或密码错误"))
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	// check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			h.errorResponse(w, r, errors.New("用户名不存在或密码错误"))
			return
		default:
			h.internalServerError(w, r, err)
			return
		}
	}

	// create jwt
	expiresAt := time.Now().Add(24 * time.Hour) // expires in one day

	claims := jwt.RegisteredClaims{
		Subject:   user.Username,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		h.internalServerError(w, r, err)
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
	http.SetCookie(w, cookie)

	// response
	h.successResponse(w, r, "登录成功", user)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// remove the token from the http-only cookie
	cookie := &http.Cookie{
		Name:    "__ecnc_shift_manager_token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, cookie)

	// response
	h.successResponse(w, r, "登出成功", nil)
}
