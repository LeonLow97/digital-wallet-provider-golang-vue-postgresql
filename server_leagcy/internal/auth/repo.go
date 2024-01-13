package auth

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/internal/utils"
	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo) {
	return &repo{
		db: db,
	}
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
		if err == sql.ErrNoRows {
			return nil, utils.UnauthorizedError{Message: "Incorrect username/password. Please try again."}
		}
		return nil, utils.InternalServerError{Message: "[Repo] error in GetByUsername: " + err.Error()}
	}
	return &user, nil
}
