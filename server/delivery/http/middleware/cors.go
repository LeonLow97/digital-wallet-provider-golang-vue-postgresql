package middleware

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/infrastructure"
)

type CorsMiddleware struct {
	cfg infrastructure.Config
}

func NewCorsMiddleware(cfg infrastructure.Config) *CorsMiddleware {
	return &CorsMiddleware{
		cfg: cfg,
	}
}

// Middleware adds CORS headers to the response.
func (m CorsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", m.cfg.Env.FrontendURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
