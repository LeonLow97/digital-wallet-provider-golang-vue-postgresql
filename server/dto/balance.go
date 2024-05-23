package dto

import "strings"

type DepositRequest struct {
	Balance  float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required"`
	UserID   int     `json:"-"`
}

type GetBalanceResponse struct {
	ID        int     `json:"id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency" validate:"required"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type GetBalanceHistory struct {
	BalanceHistory []BalanceHistory `json:"balanceHistory"`
}

type BalanceHistory struct {
	ID        int     `db:"id" json:"-"`
	UserID    int     `db:"user_id" json:"-"`
	BalanceID int     `db:"balance_id" json:"-"`
	Amount    float64 `db:"amount" json:"amount"`
	Currency  string  `db:"currency" json:"currency" validate:"required"`
	Type      string  `db:"type" json:"type"`
	CreatedAt string  `db:"created_at" json:"createdAt"`
}

type GetBalancesResponse struct {
	Balances []GetBalanceResponse `json:"balances"`
}

type WithdrawRequest struct {
	Balance  float64 `json:"amount" validate:"required,gt=0"`
	Currency string  `json:"currency" validate:"required"`
	UserID   int     `json:"-"`
}

func (req *DepositRequest) DepositSanitize() {
	req.Currency = strings.TrimSpace(req.Currency)
}

func (req *WithdrawRequest) WithdrawSanitize() {
	req.Currency = strings.TrimSpace(req.Currency)
}
