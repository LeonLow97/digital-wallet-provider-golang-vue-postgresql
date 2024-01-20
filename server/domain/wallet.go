package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Wallet struct {
	ID        int     `db:"id"`
	Type      string  `db:"type"`
	TypeID    int     `db:"type_id"`
	Balance   float64 `db:"balance"`
	Currency  string  `db:"currency"`
	CreatedAt string  `db:"created_at"`
	UserID    int     `db:"user_id"`
}

type WalletUsecase interface {
	CreateWallet(ctx context.Context, req dto.CreateWalletRequest) error
}

type WalletRepository interface {
	GetWalletTypes(ctx context.Context) (map[string]int, error)
	CheckWalletExistsByUserID(ctx context.Context, userID int, walletType string) (int, error)
	CreateWallet(ctx context.Context, wallet *Wallet) error
}
