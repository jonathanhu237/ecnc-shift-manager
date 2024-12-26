package models

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	FullName     string    `json:"fullName"`
	Role         string    `json:"role"`
	Level        int32     `json:"level"`
	CreatedAt    time.Time `json:"created_at"`
	Version      int32     `json:"-"`
}

func (m *Models) InsertUser(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, role_id)
		VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name = $5))
		RETURNING id, (SELECT level FROM roles WHERE name = $5), created_at, version
	`
	args := []any{user.Username, user.Email, user.PasswordHash, user.FullName, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Level, &user.CreatedAt, &user.Version); err != nil {
		return err
	}

	return nil
}

func (m *Models) SelectUserByUsername(username string) (*User, error) {
	user := &User{Username: username}

	query := `
		SELECT u.id, u.password_hash, u.email, u.full_name, r.name, r.level, u.created_at, u.version
		FROM users u
		INNER JOIN roles r
		ON u.role_id = r.id
		WHERE username = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Level,
		&user.CreatedAt,
		&user.Version,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Models) SelectUserByID(userID int64) (*User, error) {
	user := &User{ID: userID}

	query := `
		SELECT u.username, u.password_hash, u.email, u.full_name, r.name, r.level, u.created_at, u.version
		FROM users AS u
		INNER JOIN roles AS r
		ON u.role_id = r.id
		WHERE u.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Level,
		&user.CreatedAt,
		&user.Version,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Models) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET
			password_hash = $1,
			email = $2,
			role_id = (
				SELECT id
				FROM roles
				WHERE name = $3
			),
			version = version + 1
		WHERE id = $4 AND version = $5
	`
	args := []any{user.PasswordHash, user.Email, user.Role, user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *Models) SelectAllUsers() ([]*User, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.password_hash,
			u.email,
			u.full_name,
			r.name,
			r.level,
			u.created_at
			u.version
		FROM users AS u
		INNER JOIN roles AS r
			ON u.role_id = r.id
		ORDER BY u.created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.Email,
			&user.FullName,
			&user.Role,
			&user.Level,
			&user.CreatedAt,
			&user.Version,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *Models) DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := m.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
