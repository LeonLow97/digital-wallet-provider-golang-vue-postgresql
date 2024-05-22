package exception

import "errors"

var (
	ErrBalanceNotFound        = errors.New("balance not found")
	ErrBalanceHistoryNotFound = errors.New("balance history not found")
	ErrInsufficientFunds      = errors.New("insufficient funds for withdrawal")
)
