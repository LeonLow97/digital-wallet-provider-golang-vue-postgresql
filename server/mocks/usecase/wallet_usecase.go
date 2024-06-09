package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/stretchr/testify/mock"
)

type WalletUsecase struct {
	mock.Mock
}

type WalletUsecaseReturnValues struct {
	GetWallet      []interface{}
	GetWallets     []interface{}
	GetWalletTypes []interface{}

	CreateWallet  []interface{}
	TopUpWallet   []interface{}
	CashOutWallet []interface{}
}

func (m *WalletUsecase) GetWallet(ctx context.Context, userID, walletID int) (*domain.Wallet, error) {
	args := m.Called(ctx, userID, walletID)

	var wallet domain.Wallet
	if v, ok := args.Get(0).(domain.Wallet); ok {
		wallet = v
	}

	return &wallet, args.Error(1)
}

func (m *WalletUsecase) GetWallets(ctx context.Context, userID int) (*[]domain.Wallet, error) {
	args := m.Called(ctx, userID)

	var wallets []domain.Wallet
	if v, ok := args.Get(0).([]domain.Wallet); ok {
		wallets = v
	}

	return &wallets, args.Error(1)
}

func (m *WalletUsecase) GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error) {
	args := m.Called(ctx)

	var walletTypes []dto.GetWalletTypesResponse
	if v, ok := args.Get(0).([]dto.GetWalletTypesResponse); ok {
		walletTypes = v
	}

	return &walletTypes, args.Error(1)
}

func (m *WalletUsecase) CreateWallet(ctx context.Context, userID int, req dto.CreateWalletRequest) error {
	args := m.Called(ctx, userID, req)

	return args.Error(0)
}

func (m *WalletUsecase) TopUpWallet(ctx context.Context, userID, walletID int, req dto.UpdateWalletRequest) error {
	args := m.Called(ctx, userID, walletID, req)

	return args.Error(0)
}

func (m *WalletUsecase) CashOutWallet(ctx context.Context, userID, walletID int, req dto.UpdateWalletRequest) error {
	args := m.Called(ctx, userID, walletID, req)

	return args.Error(0)
}
