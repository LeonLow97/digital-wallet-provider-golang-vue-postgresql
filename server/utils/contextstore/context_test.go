package contextstore_test

import (
	"context"
	"testing"

	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
	"github.com/stretchr/testify/require"
)

func TestSessionIDFromContext(t *testing.T) {
	ctx := context.Background()
	sessionID := "b48c83cf-9943-4d71-b725-6e1ffd4982cc"
	ctx = contextstore.SessionIDWithContext(ctx, sessionID)

	got, err := contextstore.SessionIDFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, sessionID, got)

	ctx = context.Background()
	_, err = contextstore.SessionIDFromContext(ctx)
	require.Error(t, err, "expecting error because session ID not in context")
}

func TestUserIDFromContext(t *testing.T) {
	ctx := context.Background()
	userID := 123
	ctx = contextstore.UserIDWithContext(ctx, userID)

	got, err := contextstore.UserIDFromContext(ctx)
	require.NoError(t, err)
	require.Equal(t, userID, got)

	ctx = context.Background()
	_, err = contextstore.UserIDFromContext(ctx)
	require.Error(t, err, "expecting error because user ID not in context")
}
