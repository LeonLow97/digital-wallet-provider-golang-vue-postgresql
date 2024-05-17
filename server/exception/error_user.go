package exception

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserFound    = errors.New("user already exists")

	ErrInvalidPassword = errors.New("invalid password format")
	ErrSamePassword    = errors.New("new password cannot be the same as the current password")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is inactive")
)

var (
	ErrTOTPSecretExists   = errors.New("user totp secret already exists")
	ErrTOTPSecretNotFound = errors.New("user totp secret not found")
	ErrInvalidMFACode     = errors.New("invalid mfa code")
)
