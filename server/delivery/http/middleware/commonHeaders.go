package middleware

import (
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/utils/constants/headers"
)

type HeadersMiddleware struct{}

func NewHeadersMiddleware() HeadersMiddleware {
	return HeadersMiddleware{}
}

func (m HeadersMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add common security headers
		w.Header().Set(headers.StrictTransportSecurity, "max-age=31536000; includeSubDomains; preload") // set to 1 year, once confident on HTTPS, increase to max-age=63072000 (recommended by hstspreload)
		w.Header().Set(headers.ContentSecurityPolicy, "default-src 'self'; img-src 'none'; style-src 'self'; script-src 'none'; object-src 'none';")
		w.Header().Set(headers.CacheControl, "no-store, no-cache")
		w.Header().Set(headers.Pragma, "no-cache")
		w.Header().Set(headers.XContentTypeOptions, "nosniff")
		w.Header().Set(headers.XFrameOptions, "DENY")
		w.Header().Set(headers.ReferrerPolicy, "same-origin")

		next.ServeHTTP(w, r)
	})
}
