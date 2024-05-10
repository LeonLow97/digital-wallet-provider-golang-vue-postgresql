package utils

import "time"

const SESSION_EXPIRY = 15 * time.Minute
const PASSWORD_RESET_AUTH_TOKEN_EXPIRY = 7 * 24 * time.Hour

// Cookie Name
const JWT_COOKIE = "mw-token"

// Transaction Status
const SUBMITTED = "SUBMITTED"
const PENDING = "PENDING"
const COMPLETED = "COMPLETED"
