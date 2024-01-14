package exception

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserFound    = errors.New("user already exists")

	ErrInvalidPassword = errors.New("invalid password format")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is inactive")
)
