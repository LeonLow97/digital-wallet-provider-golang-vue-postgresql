package domain

import (
	"context"
	"database/sql"

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
	GetBalanceHistory(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceHistory, error)
	GetBalance(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceResponse, error)
	GetBalances(ctx context.Context, userID int) (*dto.GetBalancesResponse, error)
	GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error)
	Deposit(ctx context.Context, req dto.DepositRequest) error
	Withdraw(ctx context.Context, req dto.WithdrawRequest) error
	CurrencyExchange(ctx context.Context, userID int, req dto.CurrencyExchangeRequest) error
}

type BalanceRepository interface {
	CreateBalanceHistory(ctx context.Context, tx *sql.Tx, balance *Balance, depositedBalance float64, balanceType string) error
	GetBalanceHistory(ctx context.Context, userID, balanceID int) (*[]dto.BalanceHistory, error)
	GetBalances(ctx context.Context, tx *sql.Tx, userID int) ([]Balance, error)
	GetBalance(ctx context.Context, userID int, currency string) (*Balance, error)
	GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error)
	GetBalanceTx(ctx context.Context, tx *sql.Tx, userID int, currency string) (*Balance, error)
	GetBalanceById(ctx context.Context, userID int, balanceId int) (*Balance, error)
	CreateBalance(ctx context.Context, tx *sql.Tx, balance *Balance) error
	UpdateBalance(ctx context.Context, tx *sql.Tx, balance *Balance) error
	UpdateBalances(ctx context.Context, tx *sql.Tx, userID int, finalBalancesMap map[string]float64) error
}
