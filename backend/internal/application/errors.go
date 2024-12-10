package application

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, errorMessage any) {
	if err := app.writeJSON(w, status, map[string]any{"error": errorMessage}); err != nil {
		app.logError(r, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *Application) internalSeverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *Application) validateError(w http.ResponseWriter, r *http.Request, err error) {
	var errors = make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = fmt.Sprintf("needs to implement '%s'", err.Tag())
	}

	app.errorResponse(w, r, http.StatusBadRequest, errors)
}
