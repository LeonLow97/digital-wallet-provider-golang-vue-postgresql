package utils_test

import (
	"testing"

	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/stretchr/testify/require"
)

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password1@", true},
		{"password123", false},        // Missing uppercase letter
		{"PASSWORD123", false},        // Missing lowercase letter
		{"Password", false},           // Missing digit and special character
		{"Password@", false},          // Missing digit
		{"Password123", false},        // Missing special character
		{"Password1", false},          // Missing special character
		{"P@ssw0rd!", true},           // All criteria met
		{"P@ssword", false},           // Missing digit
		{"P@ssw0rd1234567890!", true}, // Longer password with all criteria met
		{"", false},                   // Empty password
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			valid := utils.IsValidPassword(tt.password)
			require.Equal(t, tt.expected, valid)
		})
	}
}
