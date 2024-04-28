package domain

import (
	"context"
	"time"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

// User is representing the User data struct
type User struct {
	ID           int     `db:"id"`
	FirstName    string  `db:"first_name"`
	LastName     string  `db:"last_name"`
	Username     string  `db:"username"`
	Email        string  `db:"email"`
	Password     string  `db:"password"`
	MobileNumber string  `db:"mobile_number"`
	Active       bool    `db:"active"`
	Admin        bool    `db:"admin"`
	Balance      float64 `db:"balance"`
}

// UserUsecase represents the user's use cases
type UserUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *dto.Token, error)
	SignUp(ctx context.Context, req dto.SignUpRequest) error
	RemoveSessionFromRedis(ctx context.Context, sessionID string) error
	GenerateJWTAccessToken(userID int, ttl time.Duration, sessionID string) (string, error)
	UpdateUser(ctx context.Context, userID int, req dto.UpdateUserRequest) error
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByEmailOrMobileNumber(ctx context.Context, email, mobileNumber string) (*User, error)
	InsertUser(ctx context.Context, user *User) error
	GetUserAndBalanceByMobileNumber(ctx context.Context, mobileNumber string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
}
