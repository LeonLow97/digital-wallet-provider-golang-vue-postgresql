package utils

import "time"

var CONSTANTS = struct {
	TIMEOUT time.Duration
}{
	TIMEOUT: 2 * time.Minute,
}

var TRANSACTION_STATUS = struct {
	COMPLETED string
	PENDING   string
	RECEIVED  string
}{
	COMPLETED: "COMPLETED",
	PENDING:   "PENDING",
	RECEIVED:  "RECEIVED",
}

type ContextUserId int

const ContextUserIdKey ContextUserId = 0
