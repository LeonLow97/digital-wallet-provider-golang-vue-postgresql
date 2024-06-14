package exception

import "errors"

var (
	ErrWalletTypeInvalid   = errors.New("wallet type invalid")
	ErrWalletAlreadyExists = errors.New("wallet already exists for this user")

	ErrNoWalletsFound = errors.New("no wallets found for this user")
	ErrNoWalletFound  = errors.New("no wallet found for this user")

	ErrInsufficientFundsForWithdrawal = errors.New("wallet has insufficient funds for withdrawal")
	ErrInsufficientFundsForDeposit    = errors.New("insufficient funds in account to top up wallet")

	ErrWalletBalanceNotFound  = errors.New("wallet balance not found")
	ErrWalletBalancesNotFound = errors.New("wallet balances not found")

	ErrWalletTypesNotFound = errors.New("wallet types not found")
)
