package dto

import "strings"

type DepositRequest struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	UserID   int     `json:"-"`
}

type BalanceResponse struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
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
