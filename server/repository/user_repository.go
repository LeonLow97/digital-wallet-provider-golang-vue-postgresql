package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	query := `
		SELECT id, email, username, password, active, admin 
		FROM users 
		WHERE email = $1;
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Active,
		&user.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmailOrMobileNumber(ctx context.Context, email, mobileNumber string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	query := `
		SELECT id, email, username, password, active, admin 
		FROM users 
		WHERE email = $1 OR mobile_number = $2;
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email, mobileNumber).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Active,
		&user.Admin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) InsertUser(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	query := `
		INSERT INTO users (first_name, last_name, email, username, password, mobile_number)
		VALUES 
		($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Username,
		user.Password,
		user.MobileNumber,
	)

	return err
}
