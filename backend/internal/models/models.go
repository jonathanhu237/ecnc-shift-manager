package models

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Users            *UserModel
	ScheduleTemplate *ScheduleTemplateModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Users:            &UserModel{DB: db},
		ScheduleTemplate: &ScheduleTemplateModel{DB: db},
	}
}
