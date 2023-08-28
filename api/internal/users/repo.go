package users

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/internal/utils"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetUserCurrencyAndBalanceByUsername(ctx context.Context, username string) (*GetUser, error)
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

func (r *repo) GetUserCurrencyAndBalanceByUsername(ctx context.Context, username string) (*GetUser, error) {
	query := `SELECT u.username, u.mobile_number, ub.balance, ub.currency
				FROM users u
				LEFT JOIN user_balance ub ON ub.user_id = u.id
				WHERE u.username = $1;`

	var Username, MobileNumber, Currency sql.NullString
	var Balance sql.NullFloat64

	query = r.db.Rebind(query)

	var user GetUser

	if err := r.db.QueryRowContext(ctx, query, username).Scan(&Username, &MobileNumber, &Balance, &Currency); err != nil {
		return nil, utils.RepositoryError{Message: err.Error()}
	}

	user = GetUser{
		Username:     Username.String,
		MobileNumber: MobileNumber.String,
		Currency:     Currency.String,
		Balance:      Balance.Float64,
	}

	return &user, nil
}
