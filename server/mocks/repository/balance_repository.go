package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type BalanceRepository struct {
	mock.Mock
}

type BalanceRepositoryReturnValues struct {
	GetBalanceHistory        []interface{}
	GetBalances              []interface{}
	GetBalance               []interface{}
	GetUserBalanceCurrencies []interface{}
	GetBalanceByID           []interface{}
	CreateBalanceHistory     []interface{}
	CreateBalance            []interface{}
	UpdateBalance            []interface{}
	UpdateBalances           []interface{}
	LogCreatorProfit         []interface{}
}

func (m *BalanceRepository) GetBalanceHistory(ctx context.Context, userID, balanceID int) (*[]dto.BalanceHistory, error) {
	args := m.Called(ctx, userID, balanceID)

	var balanceHistory *[]dto.BalanceHistory
	if v, ok := args.Get(0).(*[]dto.BalanceHistory); ok {
		balanceHistory = v
	}

	return balanceHistory, args.Error(1)
}

func (m *BalanceRepository) GetBalances(ctx context.Context, tx *sqlx.Tx, userID int) ([]domain.Balance, error) {
	args := m.Called(ctx, tx, userID)

	var balances []domain.Balance
	if v, ok := args.Get(0).([]domain.Balance); ok {
		balances = v
	}

	return balances, args.Error(1)
}

func (m *BalanceRepository) GetBalance(ctx context.Context, tx *sqlx.Tx, userID int, currency string) (*domain.Balance, error) {
	args := m.Called(ctx, tx, userID, currency)

	var balance *domain.Balance
	if v, ok := args.Get(0).(*domain.Balance); ok {
		balance = v
	}

	return balance, args.Error(1)
}

func (m *BalanceRepository) GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error) {
	args := m.Called(ctx, userID)

	var currencies *[]dto.GetUserBalanceCurrenciesResponse
	if v, ok := args.Get(0).(*[]dto.GetUserBalanceCurrenciesResponse); ok {
		currencies = v
	}

	return currencies, args.Error(1)
}

func (m *BalanceRepository) GetBalanceByID(ctx context.Context, userID int, balanceId int) (*domain.Balance, error) {
	args := m.Called(ctx, userID, balanceId)

	var balance *domain.Balance
	if v, ok := args.Get(0).(*domain.Balance); ok {
		balance = v
	}

	return balance, args.Error(1)
}

func (m *BalanceRepository) CreateBalanceHistory(ctx context.Context, tx *sqlx.Tx, balance *domain.Balance, depositedBalance float64, balanceType string) error {
	args := m.Called(ctx, tx, balance, depositedBalance, balanceType)
	return args.Error(0)
}

func (m *BalanceRepository) CreateBalance(ctx context.Context, tx *sqlx.Tx, balance *domain.Balance) error {
	args := m.Called(ctx, tx, balance)
	return args.Error(0)
}

func (m *BalanceRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, balance *domain.Balance) error {
	args := m.Called(ctx, tx, balance)
	return args.Error(0)
}

func (m *BalanceRepository) UpdateBalances(ctx context.Context, tx *sqlx.Tx, userID int, finalBalancesMap map[string]float64) error {
	args := m.Called(ctx, tx, userID, finalBalancesMap)
	return args.Error(0)
}

func (m *BalanceRepository) LogCreatorProfit(ctx context.Context, tx *sqlx.Tx, profit float64, currency string) error {
	args := m.Called(ctx, tx, profit, currency)
	return args.Error(0)
}
