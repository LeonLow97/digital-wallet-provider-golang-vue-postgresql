package middleware

import "net/http"

// SkipperFunc is a golang decorator function for skipping endpoint path
type SkipperFunc func(r *http.Request) bool

// NewSkipperFunc records the skipped endpoints and returns SkipperFunc
func NewSkipperFunc(skippedEndpoints ...string) SkipperFunc {
	skippedEndpointsMap := make(map[string]struct{})

	for _, skippedEndpoint := range skippedEndpoints {
		skippedEndpointsMap[skippedEndpoint] = struct{}{}
	}

	return func(r *http.Request) bool {
		skippedEndpoint := r.URL.EscapedPath()
		_, found := skippedEndpointsMap[skippedEndpoint]

		return found
	}
}
