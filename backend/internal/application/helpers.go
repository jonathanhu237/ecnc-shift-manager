package application

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

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

func (app *Application) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Application) successResponse(w http.ResponseWriter, message string, data any) error {
	return app.writeJSON(w, http.StatusOK, response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

func (app *Application) errorResponse(w http.ResponseWriter, code int, message string) error {
	return app.writeJSON(w, http.StatusOK, response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func (app *Application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
}

func (app *Application) internalSeverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, http.StatusInternalServerError, "internal server error")
}
