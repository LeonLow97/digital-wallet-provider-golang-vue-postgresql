package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password, active, admin
		FROM users
		WHERE username = $1`

	var user User
	row := r.db.QueryRowContext(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Active,
		&user.Admin,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
