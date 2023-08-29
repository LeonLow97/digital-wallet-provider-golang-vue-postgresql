package auth

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := User{
		Username: "validUser",
		Password: "$2a$10$jwP5FD6FyY3bG841zoQDI.c3PEBbBk17j1ZyqMjuC5NCpberWqSN.",
		Active:   1,
	}

	switch username {
	case "validUser":
		return &user, nil
	case "nonExistentUser":
		return nil, errors.New("non-existent user")
	case "inactiveUser":
		user.Active = 0
		return &user, nil
	}

	return nil, nil
}
