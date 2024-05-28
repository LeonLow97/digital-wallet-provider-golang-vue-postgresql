package dto

import (
	"strings"
)

type UpdateWalletResponse struct {
	WalletID int     `json:"wallet_id"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
}

type CreateWalletRequest struct {
	WalletTypeID   int              `json:"wallet_type_id" validate:"required,gt=0"`
	CurrencyAmount []CurrencyAmount `json:"currency_amount" validate:"required"`
}

type CurrencyAmount struct {
	Amount   float64 `json:"amount" db:"amount" validate:"required,gte=0"`
	Currency string  `json:"currency" db:"currency" validate:"required,min=0,max=3"`
}

type GetWalletsResponse struct {
	Wallets []GetWalletResponse `json:"wallets"`
}

type GetWalletResponse struct {
	WalletID       int                    `json:"wallet_id"`
	WalletTypeID   int                    `json:"wallet_type_id"`
	WalletType     string                 `json:"wallet_type"`
	CurrencyAmount []WalletCurrencyAmount `json:"wallet_currency_amount"`
	CreatedAt      string                 `db:"created_at"`
}

type WalletCurrencyAmount struct {
	WalletID  int     `json:"wallet_id" db:"id"`
	Amount    float64 `json:"amount" db:"amount"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"createdAt" db:"created_at"`
	UpdatedAt string  `json:"updatedAt" db:"updated_at"`
}

func (req *CreateWalletRequest) CreateWalletSanitize() {
	for idx, c := range req.CurrencyAmount {
		req.CurrencyAmount[idx].Currency = strings.TrimSpace(c.Currency)
	}
}

type GetWalletTypesResponse struct {
	ID         int    `json:"id" db:"id"`
	WalletType string `json:"walletType" db:"type"`
}

type WalletValidation struct {
	WalletExists      bool `db:"wallet_exists"`
	IsValidWalletType bool `db:"is_valid_wallet_type"`
}

type UpdateWalletRequest struct {
	CurrencyAmount []CurrencyAmount `json:"currency_amount" validate:"required"`
}

func (req *UpdateWalletRequest) TopUpWalletSanitize() {
	for idx, ca := range req.CurrencyAmount {
		req.CurrencyAmount[idx].Currency = strings.TrimSpace(ca.Currency)
	}
}
