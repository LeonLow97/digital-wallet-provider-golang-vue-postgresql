package apiErr

// Generic
var (
	ErrBadRequest          = "Bad Request"
	ErrInternalServerError = "Internal Server Error"
	ErrUnauthorized        = "Unauthorized"
	ErrForbidden           = "Forbidden"
)

// User
var (
	ErrInvalidCredentials = "Invalid Credentials. Please try again."
	ErrInactiveUser       = "User is inactive. Please contact the system administrator"
	ErrSuccessfulLogout   = "Logged out successfully!"

	ErrUserFound       = "Email or Mobile Number has been taken. Please try again."
	ErrInvalidPassword = `Invalid Password Format. Password must contain at least one lowercase 
						letter, one uppercase letter, one numeric digit, one special character, 
						and have a minimum length of 8 characters.`
	ErrUserNotFound             = "User does not exist. Please try again."
	ErrCurrentPasswordIncorrect = "Current password is incorrect. Please try again."
	ErrSamePassword             = "New password is same as current password. Please choose a different password."

	ErrInvalidMFACode = "Invalid MFA Code. Please try again."
)

// Wallet
var (
	ErrWalletTypeInvalid   = "Wallet type is invalid. Please try another wallet type."
	ErrWalletAlreadyExists = "The wallet you are trying to create already exist. Please try again."

	ErrNoWalletsFound = "No wallets found."
	ErrNoWalletFound  = "Wallet not found."

	ErrInsufficientFundsInAccount = "Insufficient funds in account. Please top up."

	ErrInsufficientFundsForWithdrawalFromWallet = "Unable to withdraw as the specified amount is higher than the amount in the wallet. Please reduce the amount."
	ErrInsufficientFundsForDepositToWallet      = "Unable to deposit because the account balance is lower than the specified amount. Please top up your balance."

	ErrNoWalletBalancesFound = "No wallet balances found. Please top up."
	ErrWalletTypesNotFound   = "No wallet types found."
)

// Balance
var (
	ErrBalanceNotFound = "Balance not found. Please deposit to create a new balance."

	ErrBalanceHistoryNotFound         = "Balance History not found. Please deposit to create a new balance."
	ErrInsufficientFundsForWithdrawal = "The specified amount for withdrawal is more than the current balance amount. Please lower the withdrawal amount."
	ErrWalletBalanceNotFound          = "Wallet Balance not found."

	ErrInsufficientFundsForCurrencyExchange = "Insufficient funds for currency exchange. Please deposit."

	ErrDepositCurrencyNotAllowed  = "Deposit Currency is not allowed."
	ErrWithdrawCurrencyNotAllowed = "Withdraw Currency is not allowed."

	ErrUserCurrenciesNotFound = "User currencies not found. Please inform system administrator."
)

// Beneficiary
var (
	ErrUserIDEqualBeneficiaryID = "Unable to add yourself as beneficiary. Please specify another mobile number."
	ErrBeneficiaryAlreadyExists = "Beneficiary already exists. Please specify another mobile number."

	ErrUserNotLinkedToBeneficiary    = "User is not linked to this beneficiary."
	ErrUserNotLinkedToAnyBeneficiary = "User is not linked to any beneficiary. Please add a beneficiary."
)

// Transaction
var (
	ErrUserAndWalletAssociationNotFound = "User is not associated with the specified wallet."
	ErrNoTransactionsFound              = "No transactions found."

	ErrBeneficiaryAccountNotRegistered = "Beneficiary is currently not a registered user. Please try again later or contact the System Administrator."
	ErrInsufficientFundsInWallet       = "Insufficient funds in the specified wallet. Please top up."

	ErrSenderWalletInvalid = "No wallet found with the specified Wallet ID. Please try again."
)
