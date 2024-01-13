package exception

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials. Please try again.")
	ErrInactiveUser       = errors.New("User is inactive. Please contact system administrator.")
)
