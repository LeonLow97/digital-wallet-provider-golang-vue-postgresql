package transactions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsFloat64(t *testing.T) {
	isFloat64TestCases := []struct {
		Input    interface{}
		Expected bool
	}{
		{42.0, true},           // Valid float64
		{3.14, true},           // Valid float64
		{int(42), false},       // Not a float64
		{float32(3.14), false}, // Not a float64
		{"3.14", false},        // Not a float64
	}

	for _, testCase := range isFloat64TestCases {
		t.Run(fmt.Sprintf("Input: %v, Expected: %v", testCase.Input, testCase.Expected), func(t *testing.T) {
			result := IsFloat64(testCase.Input)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

func TestValidateFloatPrecision(t *testing.T) {
	testCases := []struct {
		value          float64
		expectedError  string
		shouldSucceed  bool
	}{
		{value: 42.42, expectedError: "", shouldSucceed: true},       // Valid float with no error expected
		{value: 42.426, expectedError: "invalid float precision: 42.43", shouldSucceed: false}, // Invalid float with error expected
		{value: 0.0, expectedError: "", shouldSucceed: true},         // Valid float with no error expected
		{value: -123.456, expectedError: "invalid float precision: -123.46", shouldSucceed: false}, // Invalid negative float with error expected
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Value=%.2f", testCase.value), func(t *testing.T) {
			err := ValidateFloatPrecision(testCase.value)
			if testCase.shouldSucceed && err != nil {
				t.Errorf("Expected ValidateFloatPrecision(%v) to succeed, but it returned an error: %v", testCase.value, err)
			} else if !testCase.shouldSucceed {
				if err == nil {
					t.Errorf("Expected ValidateFloatPrecision(%v) to return an error, but it returned nil", testCase.value)
				} else if err.Error() != testCase.expectedError {
					t.Errorf("Expected error message to be %s, but got %s", testCase.expectedError, err.Error())
				}
			}
		})
	}
}
