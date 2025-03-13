package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	args := m.Called(ctx, userID)

	var user *domain.User
	if v, ok := args.Get(0).(*domain.User); ok {
		user = v
	}

	return user, args.Error(1)
}

func (m *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)

	var user *domain.User
	if v, ok := args.Get(0).(*domain.User); ok {
		user = v
	}

	return user, args.Error(1)
}

func (m *UserRepository) GetUserByEmailOrMobileNumber(ctx context.Context, email, mobileNumber string) (*domain.User, error) {
	args := m.Called(ctx, email, mobileNumber)

	var user *domain.User
	if v, ok := args.Get(0).(*domain.User); ok {
		user = v
	}

	return user, args.Error(1)
}

func (m *UserRepository) InsertUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) ChangePassword(ctx context.Context, hashedPassword string, userID int) error {
	args := m.Called(ctx, hashedPassword, userID)
	return args.Error(0)
}

func (m *UserRepository) GetUserAndBalanceByMobileNumber(ctx context.Context, mobileNumber string) (*domain.User, error) {
	args := m.Called(ctx, mobileNumber)

	var user *domain.User
	if v, ok := args.Get(0).(*domain.User); ok {
		user = v
	}

	return user, args.Error(1)
}

func (m *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) InsertUserTOTPSecret(ctx context.Context, totpConfig domain.TOTPConfiguration) error {
	args := m.Called(ctx, totpConfig)
	return args.Error(0)
}

func (m *UserRepository) UpdateIsMFAConfigured(ctx context.Context, userID int, mfaConfigured bool) error {
	args := m.Called(ctx, userID, mfaConfigured)
	return args.Error(0)
}

func (m *UserRepository) GetUserTOTPSecretCount(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)

	return args.Int(0), args.Error(1)
}

func (m *UserRepository) GetUserTOTPSecret(ctx context.Context, userID int) (string, error) {
	args := m.Called(ctx, userID)

	return args.String(0), args.Error(1)
}
