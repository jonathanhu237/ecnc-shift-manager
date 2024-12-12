package models

import (
	"context"
	"database/sql"
	"time"
)

type RefreshToken struct {
	ID               int64
	UserID           int64
	RefreshTokenHash string
	IssuedAt         time.Time
	ExpiresAt        time.Time
	Revoked          bool
}

type RefreshTokenModel struct {
	DB *sql.DB
}

func (m *RefreshTokenModel) InsertRefreshToken(rft *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, refresh_token_hash, expires_at) 
		VALUES ($1, $2, $3)
		RETURNING id, issued_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(
		ctx,
		query,
		rft.UserID,
		rft.RefreshTokenHash,
		rft.ExpiresAt,
	).Scan(&rft.ID, &rft.IssuedAt); err != nil {
		return err
	}

	return nil
}

func (m *RefreshTokenModel) RevokeUserTokens(userID int64) error {
	query := `
		UPDATE refresh_tokens 
		SET revoked = true 
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.DB.ExecContext(ctx, query, userID); err != nil {
		return err
	}

	return nil
}

func (m *RefreshTokenModel) DeleteExpiredTokens() error {
	query := `
        DELETE FROM refresh_tokens 
        WHERE expires_at < CURRENT_TIMESTAMP 
		OR revoked = true
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.DB.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
