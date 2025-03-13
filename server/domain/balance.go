package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/jmoiron/sqlx"
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
	PreviewExchange(ctx context.Context, req dto.PreviewExchangeRequest) dto.PreviewExchangeResponse
}

type BalanceRepository interface {
	GetBalanceHistory(ctx context.Context, userID, balanceID int) (*[]dto.BalanceHistory, error)
	GetBalances(ctx context.Context, tx *sqlx.Tx, userID int) ([]Balance, error)
	GetBalance(ctx context.Context, tx *sqlx.Tx, userID int, currency string) (*Balance, error)
	GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error)
	GetBalanceByID(ctx context.Context, userID int, balanceId int) (*Balance, error)

	CreateBalanceHistory(ctx context.Context, tx *sqlx.Tx, balance *Balance, depositedBalance float64, balanceType string) error
	CreateBalance(ctx context.Context, tx *sqlx.Tx, balance *Balance) error

	UpdateBalance(ctx context.Context, tx *sqlx.Tx, balance *Balance) error
	UpdateBalances(ctx context.Context, tx *sqlx.Tx, userID int, finalBalancesMap map[string]float64) error

	LogCreatorProfit(ctx context.Context, tx *sqlx.Tx, profit float64, currency string) error
}
