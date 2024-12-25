package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Application) successResponse(w http.ResponseWriter, r *http.Request, message string, data any) {
	app.writeJSON(w, r, http.StatusOK, response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeJSON(w, r, http.StatusOK, response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	})
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, err)
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
