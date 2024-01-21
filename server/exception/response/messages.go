package apiErr

// Generic
var (
	ErrBadRequest          = "Bad Request"
	ErrInternalServerError = "Internal Server Error"
	ErrUnauthorized        = "Unauthorized"
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
)

// Balance
var (
	ErrBalanceNotFound                = "Balance not found. Please deposit to create a new balance."
	ErrInsufficientFundsForWithdrawal = "The specified amount for withdrawal is more than the current balance amount. Please lower the withdrawal amount."
)
