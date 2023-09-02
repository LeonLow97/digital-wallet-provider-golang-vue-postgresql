package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGenerateJwtAccessTokenAndRefreshToken(t *testing.T) {
	testCases := []struct {
		Test              string
		User              *User
		AccessTokenExpiry time.Duration
	}{
		{
			Test: "Generate Valid Token",
			User: &User{
				ID:       1,
				Username: "validUsername",
				Active:   1,
				Admin:    1,
			},
			AccessTokenExpiry: time.Minute * 15,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// generate the token
			tokens, err := generateJwtAccessTokenAndRefreshToken(tc.User, tc.AccessTokenExpiry)
			if err != nil {
				t.Fatalf("generateJwtAccessTokenAndRefreshToken returned an error: %v", err)
			}

			// parse the access token
			accessToken, err := jwt.Parse(tokens.Token, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
			})
			if err != nil || !accessToken.Valid {
				t.Fatalf("Access token is not valid: %v", err)
			}

			// Check access token claims
			accessClaims := accessToken.Claims.(jwt.MapClaims)
			if accessClaims["name"] != tc.User.Username ||
				accessClaims["sub"] != float64(tc.User.ID) ||
				accessClaims["aud"] != issuer ||
				accessClaims["iss"] != issuer ||
				accessClaims["admin"] != true {
				t.Fatalf("expected access token claims to be: %v, %v, %v, %v, %v \nbut got: %v, %v, %v, %v, %v", tc.User.Username, tc.User.ID, issuer, issuer, true, accessClaims["name"], accessClaims["sub"], accessClaims["aud"], accessClaims["iss"], accessClaims["admin"])
			}

			// Check access token expiration
			accessExp := time.Unix(int64(accessClaims["exp"].(float64)), 0)
			if time.Until(accessExp) > tc.AccessTokenExpiry {
				t.Fatalf("Access token expiration is incorrect")
			}

			// Parse refresh token
			refreshToken, err := jwt.Parse(tokens.RefreshToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
			})
			if err != nil || !refreshToken.Valid {
				t.Fatalf("Refresh token is not valid: %v", err)
			}

			// Check refresh token claims
			refreshClaims := refreshToken.Claims.(jwt.MapClaims)
			if refreshClaims["sub"] != float64(tc.User.ID) {
				t.Fatalf("expected refresh token claims to have id: %v, but got %v", refreshClaims["sub"], float64(tc.User.ID))
			}

			// Check refresh token expiration
			refreshExp := time.Unix(int64(refreshClaims["exp"].(float64)), 0)
			if time.Until(refreshExp) > refreshTokenExpiry {
				t.Fatalf("Refresh token expiration is incorrect")
			}
		})
	}
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
