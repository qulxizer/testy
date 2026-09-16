package db

import (
	"context"
	"time"
)

func (d *DB) CreateSession(ctx context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3);
	`
	_, err := d.Pool.Exec(ctx, query, tokenHash, userID, expiresAt)
	return err
}

func (d *DB) GetUserIDBySession(ctx context.Context, tokenHash []byte) (int64, error) {
	var userID int64
	query := `
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1 AND expires_at > NOW();
	`
	err := d.Pool.QueryRow(ctx, query, tokenHash).Scan(&userID)
	return userID, err
}

// DeleteSession kills a session immediately on logout.
func (d *DB) DeleteSession(ctx context.Context, tokenHash []byte) error {
	query := `
		DELETE FROM sessions
		WHERE token_hash = $1;
	`
	_, err := d.Pool.Exec(ctx, query, tokenHash)
	return err
}
