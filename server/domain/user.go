package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

// User is representing the User data struct
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Active   bool   `db:"active"`
	Admin    bool   `db:"admin"`
}

// UserUsecase represents the user's use cases
type UserUsecase interface {
	Login(ctx context.Context, dto dto.LoginRequest) (*User, *dto.Token, error)
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
