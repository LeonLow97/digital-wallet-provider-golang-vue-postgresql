package context

import "errors"

var (
	ErrUserIDNotInContext    = errors.New("UserID not found in context")
	ErrSessionIDNotInContext = errors.New("SessionID not found in context")
)
