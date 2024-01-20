package dto

import "strings"

type CreateWalletRequest struct {
	Type     string  `json:"type" validate:"required,min=5,max=100"`
	Balance  float64 `json:"balance" validate:"required,gte=0"`
	Currency string  `json:"currency" validate:"required,min=0,max=3"`
	UserID   int     `json:"-"`
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
	req.Currency = strings.TrimSpace(req.Currency)
}
