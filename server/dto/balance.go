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

type GetUserBalanceCurrenciesResponse struct {
	Currency string `json:"currency" db:"currency"`
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

type CurrencyExchangeRequest struct {
	FromAmount float64 `json:"from_amount" validate:"required,gt=0"`
	ToCurrency string  `json:"to_currency" validate:"required,len=3"`
}

func (req *CurrencyExchangeRequest) CurrencyExchangeSanitize() {
	req.ToCurrency = strings.TrimSpace(req.ToCurrency)
}

type PreviewExchangeRequest struct {
	ActionType   string  `json:"action_type" validate:"oneof=amountToSend amountToReceive"`
	FromAmount   float64 `json:"from_amount" validate:"omitempty,gt=0"`
	FromCurrency string  `json:"from_currency" validate:"omitempty,len=3"`
	ToAmount     float64 `json:"to_amount" validate:"omitempty,gt=0"`
	ToCurrency   string  `json:"to_currency" validate:"omitempty,len=3"`
}

type PreviewExchangeResponse struct {
	ActionType   string  `json:"actionType"`
	FromAmount   float64 `json:"fromAmount"`
	FromCurrency string  `json:"fromCurrency"`
	ToAmount     float64 `json:"toAmount"`
	ToCurrency   string  `json:"toCurrency"`
}

func (req *PreviewExchangeRequest) PreviewExchangeSanitize() {
	req.FromCurrency = strings.TrimSpace(req.FromCurrency)
	req.ToCurrency = strings.TrimSpace(req.ToCurrency)
}
