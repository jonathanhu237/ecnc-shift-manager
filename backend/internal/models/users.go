package models

import "database/sql"

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	RoleID       int64  `json:"role_id"`
}

type UserModel struct {
	DB *sql.DB
}
