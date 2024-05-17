package middleware

import (
	"net/http"
	"strings"
)

// NewCSPMiddleware creates a middleware that sets Content Security Policy (CSP) headers
// to enhance security by restricting sources for content such as scripts and styles.
// This helps in mitigating risks such as cross-site scripting (XSS).
func NewCSPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			policies := []string{
				"default-src: 'self';",
				"script-src 'self';",
				"style-src 'self';",
			}

			cspPolicies := strings.Join(policies, " ")

			w.Header().Set("Content-Security-Policy", cspPolicies)

			next.ServeHTTP(w, r)
		})
	}
}
