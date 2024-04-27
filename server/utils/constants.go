package utils

import "time"

const SESSION_EXPIRY = 15 * time.Minute

// Cookie Name
const JWT_COOKIE = "mw-token"

// Transaction Status
const SUBMITTED = "SUBMITTED"
const PENDING = "PENDING"
const COMPLETED = "COMPLETED"
