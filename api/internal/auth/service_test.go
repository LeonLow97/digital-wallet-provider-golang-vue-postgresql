package auth

import (
	"context"
	"testing"
)

func Test_LoginService(t *testing.T) {
	var creds Credentials

	creds.Username = "validUser"
	creds.Password = "password"

	mockRepo := &MockRepo{}
	s := &service{
		repo: mockRepo,
	}

	// Testing valid username and password
	user, _, err := s.Login(context.Background(), &creds)
	if err != nil {
		t.Errorf("valid username and password but got error, %q", err.Error())
	}
	if user.Username != creds.Username {
		t.Errorf("expected username and password to be %s and %s, but got %s and %s", creds.Username, creds.Password, user.Username, user.Password)
	}

	// Testing invalid password
	creds.Password = "wrongpassword"
	_, _, err = s.Login(context.Background(), &creds)
	expectedError := "Incorrect username/password. Please try again."
	if err == nil || err.Error() != expectedError {
		t.Errorf("\nexpected error message: %s \nbut got: %s", expectedError, err)
	}

	// Non existent user
	creds.Username = "nonExistentUser"
	creds.Password = "password"
	_, _, err = s.Login(context.Background(), &creds)
	expectedError = "Incorrect username/password. Please try again."
	if err == nil || err.Error() != expectedError {
		t.Errorf("\nexpected error message: %s \nbut got: %s", expectedError, err)
	}

	// Inactive user
	creds.Username = "inactiveUser"
	_, _, err = s.Login(context.Background(), &creds)
	expectedError = "This user account has been disabled. Please contact the system administrator."
	if err == nil || err.Error() != expectedError {
		t.Errorf("\nexpected error message: %s \nbut got: %s", expectedError, err)
	}
}
