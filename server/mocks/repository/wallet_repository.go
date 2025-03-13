package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type WalletRepository struct {
	mock.Mock
}

type WalletRepositoryReturnValues struct {
	GetWalletByWalletID                  []interface{}
	GetWallets                           []interface{}
	GetWalletTypes                       []interface{}
	GetWalletBalancesByUserID            []interface{}
	GetWalletBalancesByUserIDAndWalletID []interface{}
	CheckWalletExistsByWalletTypeID      []interface{}
	CheckWalletTypeExists                []interface{}
	CreateWallet                         []interface{}
	InsertWalletCurrencyAmount           []interface{}
	TopUpWalletBalances                  []interface{}
	CashOutWalletBalances                []interface{}
}

func (m *WalletRepository) GetWalletByWalletID(ctx context.Context, userID, walletID int) (*domain.Wallet, error) {
	args := m.Called(ctx, userID, walletID)

	var wallet *domain.Wallet
	if v, ok := args.Get(0).(*domain.Wallet); ok {
		wallet = v
	}

	return wallet, args.Error(1)
}

func (m *WalletRepository) GetWallets(ctx context.Context, userID int) ([]domain.Wallet, error) {
	args := m.Called(ctx, userID)

	var wallets []domain.Wallet
	if v, ok := args.Get(0).([]domain.Wallet); ok {
		wallets = v
	}

	return wallets, args.Error(1)
}

func (m *WalletRepository) GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error) {
	args := m.Called(ctx)

	var getWalletTypesResponse *[]dto.GetWalletTypesResponse
	if v, ok := args.Get(0).(*[]dto.GetWalletTypesResponse); ok {
		getWalletTypesResponse = v
	}

	return getWalletTypesResponse, args.Error(1)
}

func (m *WalletRepository) GetWalletBalancesByUserID(ctx context.Context, userID int) ([]domain.WalletCurrencyAmount, error) {
	args := m.Called(ctx, userID)

	var wallets []domain.WalletCurrencyAmount
	if v, ok := args.Get(0).([]domain.WalletCurrencyAmount); ok {
		wallets = v
	}

	return wallets, args.Error(1)
}

func (m *WalletRepository) GetWalletBalancesByUserIDAndWalletID(ctx context.Context, tx *sqlx.Tx, userID, walletID int) ([]domain.WalletCurrencyAmount, error) {
	args := m.Called(ctx, tx, userID, walletID)

	var walletCurrencyAmount []domain.WalletCurrencyAmount
	if v, ok := args.Get(0).([]domain.WalletCurrencyAmount); ok {
		walletCurrencyAmount = v
	}

	return walletCurrencyAmount, args.Error(1)
}

func (m *WalletRepository) CheckWalletExistsByWalletTypeID(ctx context.Context, userID, walletTypeID int) (bool, error) {
	args := m.Called(ctx, userID, walletTypeID)
	return args.Bool(0), args.Error(1)
}

func (m *WalletRepository) CheckWalletTypeExists(ctx context.Context, WalletTypeID int) (bool, error) {
	args := m.Called(ctx, WalletTypeID)
	return args.Bool(0), args.Error(1)
}

func (m *WalletRepository) CreateWallet(ctx context.Context, tx *sqlx.Tx, wallet *domain.Wallet) (int, error) {
	args := m.Called(ctx, tx, wallet)
	return args.Int(0), args.Error(1)
}

func (m *WalletRepository) InsertWalletCurrencyAmount(ctx context.Context, tx *sqlx.Tx, walletID, userID int, currencyAmount []domain.WalletCurrencyAmount) error {
	args := m.Called(ctx, tx, walletID, userID, currencyAmount)
	return args.Error(0)
}

func (m *WalletRepository) TopUpWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	args := m.Called(ctx, tx, userID, walletID, finalWalletBalancesMap)
	return args.Error(0)
}

func (m *WalletRepository) CashOutWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	args := m.Called(ctx, tx, userID, walletID, finalWalletBalancesMap)
	return args.Error(0)
}
