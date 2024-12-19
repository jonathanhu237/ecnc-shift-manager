package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jonathanhu237/ecnc-shift-manager/backend/internal/config"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://postgres:%s@localhost:5432/ecnc_shift_manager_db?sslmode=disable", cfg.PostgresPassword))
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