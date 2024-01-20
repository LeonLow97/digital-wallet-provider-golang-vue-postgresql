package exception

import "errors"

var (
	ErrWalletTypeInvalid   = errors.New("wallet type invalid")
	ErrWalletAlreadyExists = errors.New("wallet already exists for this user")
)
