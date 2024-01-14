package apiErr

var (
	ErrBadRequest          = "Bad Request"
	ErrInternalServerError = "Internal Server Error"
)

var (
	ErrInvalidCredentials = "Invalid Credentials. Please try again."
	ErrInactiveUser       = "User is inactive. Please contact the system administrator"
	ErrSuccessfulLogout   = "Logged out successfully!"

	ErrUserFound       = "Email or Mobile Number has been taken. Please try again."
	ErrInvalidPassword = `Invalid Password Format. Password must contain at least one lowercase 
						letter, one uppercase letter, one numeric digit, one special character, 
						and have a minimum length of 8 characters.`
)
