package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/stretchr/testify/mock"
)

type BalanceUsecase struct {
	mock.Mock
}

type BalanceUsecaseReturnValues struct {
	GetBalanceHistory        []interface{}
	GetBalance               []interface{}
	GetBalances              []interface{}
	GetUserBalanceCurrencies []interface{}
	Deposit                  []interface{}
	Withdraw                 []interface{}
	CurrencyExchange         []interface{}
	PreviewExchange          interface{}
}

func (m *BalanceUsecase) GetBalanceHistory(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceHistory, error) {
	args := m.Called(ctx, userID, balanceID)

	var getBalanceHistory *dto.GetBalanceHistory
	if v, ok := args.Get(0).(*dto.GetBalanceHistory); ok {
		getBalanceHistory = v
	}

	return getBalanceHistory, args.Error(1)
}

func (m *BalanceUsecase) GetBalance(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceResponse, error) {
	args := m.Called(ctx, userID, balanceID)

	var getBalanceResponse *dto.GetBalanceResponse
	if v, ok := args.Get(0).(*dto.GetBalanceResponse); ok {
		getBalanceResponse = v
	}

	return getBalanceResponse, args.Error(1)
}

func (m *BalanceUsecase) GetBalances(ctx context.Context, userID int) (*dto.GetBalancesResponse, error) {
	args := m.Called(ctx, userID)

	var getBalancesResponse *dto.GetBalancesResponse
	if v, ok := args.Get(0).(*dto.GetBalancesResponse); ok {
		getBalancesResponse = v
	}
	return getBalancesResponse, args.Error(1)
}

func (m *BalanceUsecase) GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error) {
	args := m.Called(ctx, userID)
	if v, ok := args.Get(0).(*[]dto.GetUserBalanceCurrenciesResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *BalanceUsecase) Deposit(ctx context.Context, req dto.DepositRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *BalanceUsecase) Withdraw(ctx context.Context, req dto.WithdrawRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *BalanceUsecase) CurrencyExchange(ctx context.Context, userID int, req dto.CurrencyExchangeRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *BalanceUsecase) PreviewExchange(ctx context.Context, req dto.PreviewExchangeRequest) dto.PreviewExchangeResponse {
	args := m.Called(ctx, req)

	var previewExchangeResponse dto.PreviewExchangeResponse
	if v, ok := args.Get(0).(dto.PreviewExchangeResponse); ok {
		previewExchangeResponse = v
	}

	return previewExchangeResponse
}
