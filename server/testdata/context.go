package testdata

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
)

func InjectUserIDIntoContext(ctx context.Context, userID int) context.Context {
	if userID == 0 {
		return context.Background()
	}
	ctx = context.WithValue(ctx, contextstore.UserIDKey, userID)
	return ctx
}
