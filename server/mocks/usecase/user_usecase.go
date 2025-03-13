package mocks

import (
	"context"
	"time"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/stretchr/testify/mock"
)

type UserUsecase struct {
	mock.Mock
}

type UserUsecaseReturnValues struct {
	Login                    []interface{}
	SignUp                   []interface{}
	ChangePassword           []interface{}
	RemoveSessionFromRedis   []interface{}
	GenerateJWTAccessToken   []interface{}
	UpdateUser               []interface{}
	ExtendUserSessionInRedis []interface{}
	SendPasswordResetEmail   []interface{}
	PasswordReset            []interface{}
	ConfigureMFA             []interface{}
	VerifyMFA                []interface{}
}

func (m *UserUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)

	var loginResponse *dto.LoginResponse
	if v, ok := args.Get(0).(*dto.LoginResponse); ok {
		loginResponse = v
	}

	return loginResponse, args.Error(1)
}

func (m *UserUsecase) SignUp(ctx context.Context, req dto.SignUpRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *UserUsecase) ChangePassword(ctx context.Context, userID int, req dto.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *UserUsecase) RemoveSessionFromRedis(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *UserUsecase) GenerateJWTAccessToken(userID int, ttl time.Duration, sessionID string) (string, error) {
	args := m.Called(userID, ttl, sessionID)

	return args.String(0), args.Error(1)
}

func (m *UserUsecase) UpdateUser(ctx context.Context, userID int, req dto.UpdateUserRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *UserUsecase) ExtendUserSessionInRedis(ctx context.Context, sessionID string, sessionExpiryInMinutes time.Duration) (string, error) {
	args := m.Called(ctx, sessionID, sessionExpiryInMinutes)

	return args.String(0), args.Error(1)
}

func (m *UserUsecase) SendPasswordResetEmail(ctx context.Context, req dto.SendPasswordResetEmailRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *UserUsecase) PasswordReset(ctx context.Context, req dto.PasswordResetRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *UserUsecase) ConfigureMFA(ctx context.Context, req dto.ConfigureMFARequest) (*dto.Token, error) {
	args := m.Called(ctx, req)

	var token *dto.Token
	if v, ok := args.Get(0).(*dto.Token); ok {
		token = v
	}

	return token, args.Error(1)
}

func (m *UserUsecase) VerifyMFA(ctx context.Context, req dto.VerifyMFARequest) (*dto.Token, error) {
	args := m.Called(ctx, req)

	var token *dto.Token
	if v, ok := args.Get(0).(*dto.Token); ok {
		token = v
	}

	return token, args.Error(1)
}
