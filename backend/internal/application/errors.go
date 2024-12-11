package application

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, errorMessage string) {
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
	errs := err.(validator.ValidationErrors)
	firstError := errs[0]
	app.errorResponse(w, r, http.StatusBadRequest, fmt.Sprintf("validation for '%s' failed on '%s' tag", firstError.Field(), firstError.Tag()))
}
