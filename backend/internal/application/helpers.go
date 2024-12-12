package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (app *Application) openDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", app.config.Database.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}

func (app *Application) readJSON(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}
	return nil
}

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *Application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logError(r, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Application) successResponse(w http.ResponseWriter, r *http.Request, message string, data any) {
	app.writeJSON(w, r, http.StatusOK, response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, err appError) {
	app.writeJSON(w, r, http.StatusOK, response{
		Code:    err.code,
		Message: err.message,
		Data:    nil,
	})
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, appError{code: http.StatusBadRequest, message: err.Error()})
}

func (app *Application) internalSeverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, errInternalServer)
}

func (app *Application) validateError(w http.ResponseWriter, r *http.Request, err error) {
	errs := err.(validator.ValidationErrors)
	message := fmt.Sprintf("validator for '%s' failed on the '%s'", errs[0].Field(), errs[0].Tag())
	app.badRequest(w, r, errors.New(message))
}

func (app *Application) generateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func (app *Application) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
