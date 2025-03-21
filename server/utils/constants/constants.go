package constants

import "time"

const SESSION_EXPIRY = 15 * time.Minute
const PASSWORD_RESET_AUTH_TOKEN_EXPIRY = 7 * 24 * time.Hour

// Cookie Name
const JWT_COOKIE = "mw-token"

// Transaction Status
const (
	FAILED  = "FAILED"
	PENDING = "PENDING"
	SUCCESS = "SUCCESS"
)

// Assumption: Depending on the user's mobile country code, the user
// will be allowed to deposit and withdraw with that currency
var CountryCodeToCurrencyMap = map[string]string{
	"+65": "SGD", "+60": "MYR", "+1": "AUD", "+61": "USD",
}

// Allowed toCurrencies, TODO: store this in database instead
var ToCurrencies = map[string]struct{}{
	"SGD": {}, "MYR": {}, "AUD": {}, "USD": {},
}

// Pagination limits
const (
	DEFAULT_PAGE_SIZE = 10
	MAX_PAGE_SIZE     = 100
)
