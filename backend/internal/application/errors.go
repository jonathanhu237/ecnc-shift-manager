package application

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, error any) {
	if err := app.writeJSON(w, status, map[string]any{"error": error}); err != nil {
		app.logError(r, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *Application) internalSeverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		errors := make(map[string]string)
		for _, err := range validationErrors {
			errors[err.Field()] = fmt.Sprintf("validation failed on '%s' tag", err.Tag())
		}
		app.errorResponse(w, r, http.StatusBadRequest, errors)
		return
	}
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusNotFound, err.Error())
}

func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}
