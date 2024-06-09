package testdata

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/utils"
)

func InjectUserIDIntoContext(ctx context.Context, userID int) context.Context {
	if userID == 0 {
		return context.Background()
	}
	ctx = context.WithValue(ctx, utils.UserIDKey, userID)
	return ctx
}
