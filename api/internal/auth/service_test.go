package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mocking the repository layer
type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*User), args.Error(1)
}

func Test_Login_Service(t *testing.T) {
	testCases := []struct {
		Test         string
		Credentials  *Credentials
		ExpectedUser *User
		ExpectErr    bool
		MockError    error
	}{
		{
			Test: "Successful Login",
			Credentials: &Credentials{
				Username: "validUsername",
				Password: "validPassword",
			},
			ExpectedUser: &User{
				ID:       1,
				Username: "validUsername",
				Password: "$2a$10$Tre6ATlvkBUWdrdR1fA.w.3d9CfJ86eHbcpyHSeAhrZaHTpmDCmIy",
				Active:   1,
				Admin:    1,
			},
			ExpectErr: false,
			MockError: nil,
		},
		{
			Test: "Non-Existent User",
			Credentials: &Credentials{
				Username: "invalidUsername",
				Password: "invalidPassword",
			},
			ExpectedUser: nil,
			ExpectErr:    true,
			MockError:    errors.New("this user does not exist"),
		},
		{
			Test: "Incorrect Password",
			Credentials: &Credentials{
				Username: "validUsername",
				Password: "invalidPassword",
			},
			ExpectedUser: &User{
				ID:       1,
				Username: "validUsername",
				Password: "$2a$10$Tre6ATlvkBUWdrdR1fA.w.3d9CfJ86eHbcpyHSeAhrZaHTpmDCmIy",
				Active:   1,
				Admin:    1,
			},
			ExpectErr: true,
			MockError: nil,
		},
		{
			Test: "Inactive User",
			Credentials: &Credentials{
				Username: "inactiveUser",
				Password: "validPassword",
			},
			ExpectedUser: &User{
				ID:       1,
				Username: "inactiveUser",
				Password: "$2a$10$Tre6ATlvkBUWdrdR1fA.w.3d9CfJ86eHbcpyHSeAhrZaHTpmDCmIy",
				Active:   0,
				Admin:    1,
			},
			ExpectErr: true,
			MockError: nil,
		},
	}

	// Creating a mock repository and a service with the mock repository
	mockRepo := mockRepo{}
	s, err := NewService(&mockRepo)
	require.NoError(t, err, "getting service with mock repo in LoginService")

	for _, tc := range testCases {
		t.Run(tc.Test, func(t *testing.T) {
			mockRepo.On("GetByUsername", mock.Anything, tc.Credentials.Username).Return(tc.ExpectedUser, tc.MockError)

			// calling the Login service
			user, _, err := s.Login(context.Background(), tc.Credentials)

			if !tc.ExpectErr {
				require.NoError(t, err, "expected no error")
				assert.Equal(t, user, tc.ExpectedUser)
			} else {
				require.Error(t, err, "expected error")
				assert.Nil(t, user)
			}
		})
	}
}
