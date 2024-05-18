package dto

import "strings"

type DepositRequest struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	UserID   int     `json:"-"`
}

type GetBalanceResponse struct {
	ID        int     `json:"id"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type GetBalancesResponse struct {
	Balances []GetBalanceResponse `json:"balances"`
}

type WithdrawRequest struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	UserID   int     `json:"-"`
}

func (req *DepositRequest) DepositSanitize() {
	req.Currency = strings.TrimSpace(req.Currency)
}

func (req *WithdrawRequest) WithdrawSanitize() {
	req.Currency = strings.TrimSpace(req.Currency)
}
