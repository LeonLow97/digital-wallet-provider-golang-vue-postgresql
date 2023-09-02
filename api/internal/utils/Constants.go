package utils

import "time"

var CONSTANTS = struct {
	TIMEOUT time.Duration
}{
	TIMEOUT: 2 * time.Minute,
}

type ContextUserId int

const ContextUserIdKey ContextUserId = 0
