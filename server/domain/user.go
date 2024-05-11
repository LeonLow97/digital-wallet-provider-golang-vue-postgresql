package domain

import (
	"context"
	"time"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type TOTPConfiguration struct {
	ID                  int       `json:"id" db:"id"`
	UserID              int       `json:"-"`
	Email               string    `json:"email" db:"email"`
	TOTPEncryptedSecret string    `json:"-" db:"totp_encrypted_secret"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// User is representing the User data struct
type User struct {
	ID              int     `db:"id"`
	FirstName       string  `db:"first_name"`
	LastName        string  `db:"last_name"`
	Username        string  `db:"username"`
	Email           string  `db:"email"`
	Password        string  `db:"password"`
	MobileNumber    string  `db:"mobile_number"`
	IsMFAConfigured bool    `db:"is_mfa_configured"`
	Active          bool    `db:"active"`
	Admin           bool    `db:"admin"`
	Balance         float64 `db:"balance"`
}

// UserUsecase represents the user's use cases
type UserUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *dto.Token, error)
	SignUp(ctx context.Context, req dto.SignUpRequest) error
	ChangePassword(ctx context.Context, userID int, req dto.ChangePasswordRequest) error
	RemoveSessionFromRedis(ctx context.Context, sessionID string) error
	GenerateJWTAccessToken(userID int, ttl time.Duration, sessionID string) (string, error)
	UpdateUser(ctx context.Context, userID int, req dto.UpdateUserRequest) error
	ExtendUserSessionInRedis(ctx context.Context, sessionID string, sessionExpiryInMinutes time.Duration) error
	SendPasswordResetEmail(ctx context.Context, req dto.SendPasswordResetEmailRequest) error
	PasswordReset(ctx context.Context, req dto.PasswordResetRequest) error
}

// UserRepository represents the user's repository contract
type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByEmailOrMobileNumber(ctx context.Context, email, mobileNumber string) (*User, error)
	InsertUser(ctx context.Context, user *User) error
	ChangePassword(ctx context.Context, hashedPassword string, userID int) error
	GetUserAndBalanceByMobileNumber(ctx context.Context, mobileNumber string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	InsertUserTOTPSecret(ctx context.Context, totpConfig TOTPConfiguration) error
}
