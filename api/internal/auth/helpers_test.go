package auth

import (
	"errors"
	"testing"
)

func Test_generateJwtAccessTokenAndRefreshToken(t *testing.T) {

}

func Test_passwordMatchers(t *testing.T) {
	testCases := []struct {
		Test           string
		HashedPassword string
		PlainText      string
		ExpectedResult bool
		ExpectedError  error
	}{
		{
			Test:           "Matching Passwords",
			HashedPassword: "$2a$10$qPbCijPlApG1UWkPGG5pbeoVc7tXyulslfwQRXH5yS3U2ovXH5u3e",
			PlainText:      "password123",
			ExpectedResult: true,
			ExpectedError:  nil,
		},
		{
			Test:           "Non-Matching Passwords",
			HashedPassword: "$2a$10$qPbCijPlApG1UWkPGG5pbeoVc7tXyulslfwQRXH5yS3U2ovXH5u3e",
			PlainText:      "incorrectpassword5678",
			ExpectedResult: false,
			ExpectedError:  nil,
		},
		{
			Test:           "Missing Hashed Password",
			HashedPassword: "",
			PlainText:      "incorrectpassword5678",
			ExpectedResult: false,
			ExpectedError:  errors.New("internal server error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// calling the helper method passwordMatchers
			isMatchingPassword, err := passwordMatchers(tc.HashedPassword, tc.PlainText)

			if isMatchingPassword != tc.ExpectedResult {
				t.Errorf("expected passwordMatch to return %v but got %v", tc.ExpectedResult, isMatchingPassword)
			}


			if tc.ExpectedError != nil {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			}
		})
	}
}
