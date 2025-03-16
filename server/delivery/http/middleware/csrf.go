package middleware

import (
	"log"
	"net/http"

	apiErr "github.com/LeonLow97/go-clean-architecture/exception/response"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/utils/contextstore"
	"github.com/LeonLow97/go-clean-architecture/utils/jsonutil"
)

// SkippedCSRFEndpoints is a set of paths to bypass CSRF token verification
var SkippedCSRFEndpoints = map[string]struct{}{
	"/api/v1/login":                {},
	"/api/v1/password-reset/send":  {},
	"/api/v1/password-reset/reset": {},
	"/api/v1/signup":               {},
	"/api/v1/health":               {},
	"/api/v1/configure-mfa":        {},
	"/api/v1/verify-mfa":           {},
}

type CSRFMiddleware struct {
	cfg         infrastructure.Config
	redisClient infrastructure.RedisClient
}

func NewCSRFMiddleware(cfg infrastructure.Config, redisClient infrastructure.RedisClient) *CSRFMiddleware {
	return &CSRFMiddleware{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func (m CSRFMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, exists := SkippedCSRFEndpoints[r.URL.Path]; exists {
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
