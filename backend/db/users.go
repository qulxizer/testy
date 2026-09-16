package db

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}

func (d *DB) CreateUser(ctx context.Context, uname string) (User, error) {
	var user User

	query := `
		WITH inserted AS (
			INSERT INTO users (username)
			VALUES ($1)
			ON CONFLICT (username) DO NOTHING
			RETURNING id, username
		)
		SELECT id, username FROM inserted
		UNION ALL
		SELECT id, username FROM users WHERE username = $1
		LIMIT 1;
	`

	err := d.Pool.QueryRow(ctx, query, uname).Scan(
		&user.ID,
		&user.Username,
	)

	return user, err
}

func (d *DB) GetUserByUsername(ctx context.Context, uname string) (*User, error) {
	var u User
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE username = $1;
	`
	err := d.Pool.QueryRow(ctx, query, uname).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var u User
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE id = $1;
	`
	err := d.Pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
