package middleware

import (
	"log"
	"net/http"

	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
)

type CSRFMiddleware struct {
	cfg         infrastructure.Config
	skipperFunc SkipperFunc
	redisClient infrastructure.RedisClient
}

func NewCSRFMiddleware(cfg infrastructure.Config, skipperFunc SkipperFunc, redisClient infrastructure.RedisClient) *CSRFMiddleware {
	return &CSRFMiddleware{
		cfg:         cfg,
		skipperFunc: skipperFunc,
		redisClient: redisClient,
	}
}

func (m CSRFMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (m.skipperFunc != nil && m.skipperFunc(r)) ||
			(r.Method != http.MethodPost && r.Method != http.MethodPut &&
				r.Method != http.MethodPatch && r.Method != http.MethodDelete) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		reqCsrfToken := r.Header.Get("X-CSRF-Token")

		if len(reqCsrfToken) == 0 {
			log.Println("client csrf token is empty")
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// retrieve sessionID from context
		sessionID, _ := contextstore.SessionIDFromContext(ctx)

		// retrieve csrfToken from redis client
		serverCsrfToken, err := m.redisClient.HGet(ctx, sessionID, "csrfToken")
		if err != nil {
			log.Println("failed to retrieve server csrf token with HGet redis client", err)
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		if len(serverCsrfToken) == 0 {
			log.Println("server csrf token is empty")
			jsonutil.ErrorJSON(w, apiErr.ErrInternalServerError, http.StatusInternalServerError)
			return
		}

		// compare server csrf token with client csrf token
		if reqCsrfToken != serverCsrfToken {
			log.Println("client csrf token is different from server csrf token")
			jsonutil.ErrorJSON(w, apiErr.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
