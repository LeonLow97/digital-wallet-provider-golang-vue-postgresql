package exception

import "errors"

var (
	ErrBalanceNotFound = errors.New("balance not found")
	ErrInsufficientFunds = errors.New("insufficient funds for withdrawal")
)
