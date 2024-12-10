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

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}

	return nil
}
