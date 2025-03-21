package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/jmoiron/sqlx"
)

type Wallet struct {
	ID             int                    `json:"id" db:"id"`
	WalletType     string                 `json:"walletType" db:"wallet_type"`
	WalletTypeID   int                    `json:"walletTypeID" db:"wallet_type_id"`
	UserID         int                    `json:"userID" db:"user_id"`
	CreatedAt      string                 `json:"createdAt" db:"created_at"`
	CurrencyAmount []WalletCurrencyAmount `json:"currencyAmount"`
}

type WalletCurrencyAmount struct {
	WalletID  int     `json:"wallet_id" db:"wallet_id"`
	Amount    float64 `json:"amount" db:"amount"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
	UpdatedAt string  `json:"updatedAt" db:"updated_at"`
}

type WalletUsecase interface {
	GetWallet(ctx context.Context, userID, walletID int) (*Wallet, error)
	GetWallets(ctx context.Context, userID int) (*[]Wallet, error)
	GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error)

	CreateWallet(ctx context.Context, userID int, req dto.CreateWalletRequest) error

	TopUpWallet(ctx context.Context, userID, walletID int, req dto.UpdateWalletRequest) error
	CashOutWallet(ctx context.Context, userID, walletID int, req dto.UpdateWalletRequest) error
}

type WalletRepository interface {
	GetWalletByWalletID(ctx context.Context, userID, walletID int) (*Wallet, error)
	GetWallets(ctx context.Context, userID int) ([]Wallet, error)
	GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error)
	GetWalletBalancesByUserID(ctx context.Context, userID int) ([]WalletCurrencyAmount, error)
	GetWalletBalancesByUserIDAndWalletID(ctx context.Context, tx *sqlx.Tx, userID, walletID int) ([]WalletCurrencyAmount, error)

	CheckWalletExistsByWalletTypeID(ctx context.Context, userID, walletTypeID int) (bool, error)
	CheckWalletTypeExists(ctx context.Context, WalletTypeID int) (bool, error)

	CreateWallet(ctx context.Context, tx *sqlx.Tx, wallet *Wallet) (int, error)
	InsertWalletCurrencyAmount(ctx context.Context, tx *sqlx.Tx, walletID, userID int, currencyAmount []WalletCurrencyAmount) error

	TopUpWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error
	CashOutWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error
}
