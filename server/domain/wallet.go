package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Wallet struct {
	ID           int     `db:"id"`
	Type         string  `db:"type"`
	TypeID       int     `db:"type_id"`
	Balance      float64 `db:"balance"`
	Currency     string  `db:"currency"`
	CreatedAt    string  `db:"created_at"`
	UserID       int     `db:"user_id"`
	FirstName    string  `db:"first_name"`
	LastName     string  `db:"last_name"`
	Username     string  `db:"username"`
	Email        string  `db:"email"`
	MobileNumber string  `db:"mobile_number"`
	Active       bool    `db:"active"`
}

type WalletUsecase interface {
	GetWallet(ctx context.Context, userID, walletID int) (*dto.GetWalletResponse, error)
	GetWallets(ctx context.Context, userID int) (*dto.GetWalletsResponse, error)
	CreateWallet(ctx context.Context, req dto.CreateWalletRequest) error
	UpdateWallet(ctx context.Context, req dto.UpdateWalletRequest) (*dto.UpdateWalletResponse, error)
}

type WalletRepository interface {
	GetWalletByWalletID(ctx context.Context, userID, walletID int) (*Wallet, error)
	GetWalletByWalletType(ctx context.Context, userID int, walletType string) (*Wallet, error)
	GetWallets(ctx context.Context, userID int) ([]Wallet, error)
	GetWalletTypes(ctx context.Context) (map[string]int, error)
	CheckWalletExistsByUserID(ctx context.Context, userID int, walletType string) (int, error)
	GetBalanceByUserID(ctx context.Context, userID int) (*Balance, error)
	CreateWallet(ctx context.Context, wallet *Wallet) error
	UpdateWallet(ctx context.Context, wallet *Wallet) error
	GetUserAndWalletByUserID(ctx context.Context, userID int, walletID int, walletCurrency string) (*Wallet, error)
}
