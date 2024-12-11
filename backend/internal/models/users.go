package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) CheckUserExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) CheckEmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) InsertUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name = $5))
		RETURNING id
	`
	args := []any{user.Username, user.Email, user.PasswordHash, user.FullName, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}
