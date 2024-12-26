package models

import (
	"database/sql"
)

type Models struct {
	db *sql.DB
}

func New(db *sql.DB) *Models {
	return &Models{db: db}
}
