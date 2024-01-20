package exception

import "errors"

var (
	ErrWalletTypeInvalid   = errors.New("wallet type invalid")
	ErrWalletAlreadyExists = errors.New("wallet already exists for this user")

	ErrNoWalletsFound = errors.New("no wallets found for this user")
	ErrNoWalletFound = errors.New("no wallet found for this user")
)
