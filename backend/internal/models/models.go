package models

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Users        *UserModel
	RefreshToken *RefreshTokenModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Users:        &UserModel{DB: db},
		RefreshToken: &RefreshTokenModel{DB: db},
	}
}
