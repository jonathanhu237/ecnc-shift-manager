package models

import "database/sql"

type Models struct {
	Users *UserModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Users: &UserModel{DB: db},
	}
}
