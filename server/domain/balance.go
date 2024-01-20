package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Balance struct {
	ID        int     `db:"id"`
	Balance   float64 `db:"balance"`
	Currency  string  `db:"currency"`
	UserID    int     `db:"user_id"`
	CreatedAt string  `db:"created_at"`
	UpdatedAt string  `db:"updated_at"`
}

type BalanceUsecase interface {
	Deposit(ctx context.Context, req dto.DepositRequest) (*dto.BalanceResponse, error)
	Withdraw(ctx context.Context, req dto.WithdrawRequest) (*dto.BalanceResponse, error)
}

type BalanceRepository interface {
	GetBalanceByUserID(ctx context.Context, userID int, currency string) (*Balance, error)
	CreateBalance(ctx context.Context, balance *Balance) error
	UpdateBalance(ctx context.Context, balance *Balance) error
}
