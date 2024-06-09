package testdata

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/gorilla/mux"
)

func RegularUserIDInjector(userID int) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userID == 0 {
				next.ServeHTTP(w, r)
				return
			}

			ctx := utils.UserIDWithContext(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
