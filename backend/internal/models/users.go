package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Level        int    `json:"level"`
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

func (m *UserModel) CheckBlackcoreExists() (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (SELECT 1 FROM users WHERE role_id = (SELECT id FROM roles WHERE name = '黑心'))
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query).Scan(&exists); err != nil {
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

func (m *UserModel) SelectUserByUsername(username string) (*User, error) {
	user := &User{Username: username}

	query := `
		SELECT u.id, u.password_hash, u.email, u.full_name, r.name, r.level
		FROM users u
		INNER JOIN roles r
		ON u.role_id = r.id
		WHERE username = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Level,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m *UserModel) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET
			password_hash = $1,
			email = $2,
			role_id = (
				SELECT id
				FROM roles
				WHERE name = $3
			)
		WHERE id = $4
	`
	args := []any{user.PasswordHash, user.Email, user.Role, user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return ErrRecordNotFound
	}

	return nil
}
