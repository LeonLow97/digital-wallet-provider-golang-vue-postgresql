package exception

import "errors"

var (
	ErrUserIsInactive        = errors.New("user is inactive")
	ErrBeneficiaryIsInactive = errors.New("beneficiary is inactive")
	ErrSenderWalletInvalid   = errors.New("sender wallet invalid")

	// transaction by wallet
	ErrUserAndWalletAssociationNotFound        = errors.New("user and wallet association not found")
	ErrBeneficiaryAndWalletAssociationNotFound = errors.New("beneficiary and balance association not found")
	ErrSenderWalletIDEqualBeneficiaryWalletID  = errors.New("sender wallet id is equal to beneficiary wallet id")
	ErrInsufficientFundsInWallet               = errors.New("insufficient funds for transfer in wallet")

	ErrNoTransactionsFound = errors.New("no transactions found")
)
