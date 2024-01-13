package exception

import "errors"

// Generic Error
var (
	ErrBadRequest          = errors.New("Bad Request")
	ErrInternalServerError = errors.New("Internal Server Error")
)
