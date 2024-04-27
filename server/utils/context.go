package utils

import "context"

type contextKey string

const UserIDKey contextKey = "UserID"
const SessionIDKey contextKey = "SessionID"

// SessionIDFromContext retrieves the session ID from the context
func SessionIDFromContext(ctx context.Context) (string, error) {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	if !ok {
		return "", ErrSessionIDNotInContext
	}
	return sessionID, nil
}

// SessionIDWithContext sets the session ID in the context
func SessionIDWithContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// UserIDFromContext retrieves the user ID from the context
func UserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0, ErrUserIDNotInContext
	}
	return userID, nil
}

// UserIDWithContext sets the user ID in the context
func UserIDWithContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}
