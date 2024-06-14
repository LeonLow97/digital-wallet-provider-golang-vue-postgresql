package exception

import "errors"

var (
	ErrBalanceNotFound  = errors.New("balance not found")
	ErrBalancesNotFound = errors.New("balances not found")

	ErrBalanceHistoryNotFound               = errors.New("balance history not found")
	ErrInsufficientFunds                    = errors.New("insufficient funds for withdrawal")
	ErrInsufficientFundsForCurrencyExchange = errors.New("insufficient funds for currency exchange")

	ErrDepositCurrencyNotAllowed  = errors.New("deposit currency not allowed")
	ErrWithdrawCurrencyNotAllowed = errors.New("withdraw currency not allowed")

	ErrToCurrencyNotAllowed        = errors.New("toCurrency is not allowed")
	ErrFromCurrencyEqualToCurrency = errors.New("fromCurrency is equal to toCurrency")

	ErrUserCurrenciesNotFound = errors.New("user currencies not found")
)
