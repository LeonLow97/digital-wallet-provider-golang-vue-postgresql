package domain

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/go-clean-architecture/dto"
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
	GetWallet(ctx context.Context, userID, walletID int) (*dto.GetWalletResponse, error)
	GetWallets(ctx context.Context, userID int) (*[]Wallet, error)
	GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error)
	CreateWallet(ctx context.Context, userID int, req dto.CreateWalletRequest) error
	UpdateWallet(ctx context.Context, req dto.UpdateWalletRequest) (*dto.UpdateWalletResponse, error)
}

type WalletRepository interface {
	GetWalletByWalletID(ctx context.Context, userID, walletID int) (*Wallet, error)
	GetWalletByWalletType(ctx context.Context, userID int, walletType string) (*Wallet, error)
	GetWallets(ctx context.Context, userID int) ([]Wallet, error)
	GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error)
	GetWalletBalancesByUserID(ctx context.Context, userID int) ([]WalletCurrencyAmount, error)
	PerformWalletValidationByUserID(ctx context.Context, userID, walletTypeID int) (*dto.WalletValidation, error)
	GetAllBalancesByUserID(ctx context.Context, userID int) ([]Balance, error)
	CreateWallet(ctx context.Context, tx *sql.Tx, wallet *Wallet) (int, error)
	InsertWalletCurrencyAmount(ctx context.Context, tx *sql.Tx, walletID, userID int, currencyAmount []WalletCurrencyAmount) error
	UpdateWallet(ctx context.Context, wallet *Wallet) error
}
