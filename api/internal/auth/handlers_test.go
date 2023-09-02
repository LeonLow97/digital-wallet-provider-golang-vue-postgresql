package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	mock.Mock
}

func (s *mockService) Login(ctx context.Context, creds *Credentials) (*User, *Token, error) {
	args := s.Called(ctx, creds)
	return args.Get(0).(*User), args.Get(1).(*Token), args.Error(2)
}

func Test_Login_Handler(t *testing.T) {
	testCases := []struct {
		Test                 string
		Body                 []byte
		Credentials          *Credentials
		MockError            error
		ExpectErr            bool
		ExpectedUser         *User
		ExpectedToken        *Token
		ExpectedJSONResponse string
		ExpectedStatusCode   int
	}{
		{
			Test: "Successful Login",
			Body: []byte(`{"username": "validUsername", "password": "validPassword"}`),
			Credentials: &Credentials{
				Username: "validUsername",
				Password: "validPassword",
			},
			MockError: nil,
			ExpectErr: false,
			ExpectedUser: &User{
				ID:       1,
				Username: "validUsername",
				Active:   1,
				Admin:    1,
			},
			ExpectedToken: &Token{
				Token:        "accessToken",
				RefreshToken: "refreshToken",
			},
			ExpectedJSONResponse: `{"token":{"access_token":"accessToken","refresh_token":"refreshToken"},"user":{"active":1,"admin":1,"id":1,"username":"validUsername"}}`,
			ExpectedStatusCode:   http.StatusOK,
		},
		{
			Test: "Unsuccessful Login",
			Body: []byte(`{"username": "wrongUsername", "password": "wrongPassword"}`),
			Credentials: &Credentials{
				Username: "wrongUsername",
				Password: "wrongPassword",
			},
			MockError:            errors.New("incorrect username/password"),
			ExpectErr:            true,
			ExpectedUser:         nil,
			ExpectedToken:        nil,
			ExpectedJSONResponse: `{"error":true,"message":"incorrect username/password"}`,
			ExpectedStatusCode:   http.StatusBadRequest,
		},
		{
			Test: "Invalid JSON Request Body",
			Body: []byte(`{"username": "wrongUsername", "password": "wrongPassword"}{"username: "invalidjson}`),
			Credentials: &Credentials{
				Username: "wrongUsername",
				Password: "wrongPassword",
			},
			MockError:            errors.New("Bad Request!"),
			ExpectErr:            true,
			ExpectedUser:         nil,
			ExpectedToken:        nil,
			ExpectedJSONResponse: `{"error":true,"message":"Bad Request!"}`,
			ExpectedStatusCode:   http.StatusBadRequest,
		},
	}

	// creating the mock service
	mockService := mockService{}
	authHandler, err := NewAuthHandler(&mockService)
	require.NoError(t, err, "getting authHandler with mockService in Login handler")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			// create a mock POST request and pass in the request body
			req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(string(tc.Body)))
			require.NoError(t, err)

			mockService.On("Login", mock.Anything, tc.Credentials).Return(tc.ExpectedUser, tc.ExpectedToken, tc.MockError)

			// create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			rr.Header().Set("Content-Type", "application/json")

			// calling the handler method
			handler := http.HandlerFunc(authHandler.Login)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.ExpectedStatusCode, rr.Code)

			// parse the jsonResponse
			var response envelope
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Error parsing JSON response: %v", err)
			}

			jsonData, _ := json.Marshal(response)

			if !tc.ExpectErr {
				assert.Equal(t, tc.ExpectedJSONResponse, string(jsonData))
			} else {
				assert.Equal(t, tc.ExpectedJSONResponse, string(jsonData))
			}
		})
	}
}
