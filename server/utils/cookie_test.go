package utils_test

import (
	"net/http/httptest"
	"testing"

	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/stretchr/testify/require"
)

// Constants for testing
const (
	TestToken = "test-token"
)

func TestIssueCookie(t *testing.T) {
	// Create a mock response writer
	w := httptest.NewRecorder()

	// Issue the cookie
	utils.IssueCookie(w, TestToken)

	// Retrieve the set cookies
	cookies := w.Result().Cookies()

	// Validate the cookie
	require.Equal(t, 1, len(cookies), "expected one cookie to be set")
	require.Equal(t, "mw-token", cookies[0].Name, "cookie name should be JWT-Cookie")
	require.Equal(t, TestToken, cookies[0].Value, "cookie value should match the test token")
	require.Equal(t, 3600, cookies[0].MaxAge, "cookie max age should be 3600 seconds")
	require.Equal(t, "/", cookies[0].Path, "cookie path should be /")
	require.Equal(t, "localhost", cookies[0].Domain, "cookie domain should be localhost")
	require.False(t, cookies[0].Secure, "cookie should not be secure")
	require.True(t, cookies[0].HttpOnly, "cookie should be HttpOnly")
}
