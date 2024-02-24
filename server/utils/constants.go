package utils

import "time"

const SESSION_EXPIRY = 15 * time.Minute

// Transaction Status
const SUBMITTED = "SUBMITTED"
const PENDING = "PENDING"
const COMPLETED = "COMPLETED"
