package dto

import "strings"

type UpdateWalletRequest struct {
	Type            string  `json:"type" validate:"required,min=5,max=100"`
	Balance         float64 `json:"balance" validate:"required"`
	BalanceCurrency string  `json:"balance_currency" validate:"required"`
	UserID          int     `json:"-"`
}

type UpdateWalletResponse struct {
	WalletID int     `json:"wallet_id"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
}

type CreateWalletRequest struct {
	Type            string  `json:"type" validate:"required,min=5,max=100"`
	Balance         float64 `json:"balance" validate:"required,gte=0"`
	WalletCurrency  string  `json:"wallet_currency" validate:"required,min=0,max=3"`
	BalanceCurrency string  `json:"balance_currency" validate:"required,min=0,max=3"`
	UserID          int     `json:"-"`
}

type GetWalletsResponse struct {
	Wallets []GetWalletResponse `json:"wallets"`
}

type GetWalletResponse struct {
	WalletID  int     `json:"wallet_id"`
	Type      string  `json:"type"`
	TypeID    int     `json:"type_id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	CreatedAt string  `db:"created_at"`
}

func (req *CreateWalletRequest) CreateWalletSanitize() {
	req.Type = strings.TrimSpace(req.Type)
	req.WalletCurrency = strings.TrimSpace(req.WalletCurrency)
	req.BalanceCurrency = strings.TrimSpace(req.BalanceCurrency)
}

func (req *UpdateWalletRequest) UpdateWalletSanitize() {
	req.Type = strings.TrimSpace(req.Type)
	req.BalanceCurrency = strings.TrimSpace(req.BalanceCurrency)
}
