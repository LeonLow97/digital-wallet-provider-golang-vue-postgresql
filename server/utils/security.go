package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateAuthenticationToken generates a random token of specified length
func GenerateAuthenticationToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
