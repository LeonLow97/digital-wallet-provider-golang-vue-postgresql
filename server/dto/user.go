package dto

import "strings"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=7,max=20"`
	Password string `json:"password" validate:"required,min=7,max=60"`
}

type LoginResponse struct {
	FirstName       string           `json:"firstName"`
	LastName        string           `json:"lastName"`
	Email           string           `json:"email"`
	Username        string           `json:"username"`
	MobileNumber    string           `json:"mobileNumber"`
	IsMFAConfigured bool             `json:"isMfaConfigured"`
	MFAConfig       MFAConfiguration `json:"mfaConfig,omitempty"`
}

type MFAConfiguration struct {
	Secret string `json:"secret,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type SignUpRequest struct {
	FirstName    *string `json:"first_name,omitempty" validate:"omitempty,min=3,max=50"`
	LastName     *string `json:"last_name,omitempty" validate:"omitempty,min=3,max=50"`
	Username     string  `json:"username" validate:"required,min=7,max=20"`
	Email        string  `json:"email" validate:"required,email,max=255"`
	Password     string  `json:"password" validate:"required,min=7,max=60"`
	MobileNumber string  `json:"mobile_number" validate:"required,min=5,max=255"`
}

func (req *LoginRequest) LoginSanitize() {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func (req *SignUpRequest) SignUpSanitize() {
	if req.FirstName != nil {
		*req.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		*req.LastName = strings.TrimSpace(*req.LastName)
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=7,max=60"`
	NewPassword     string `json:"new_password" validate:"required,min=7,max=60"`
}

func (req *ChangePasswordRequest) ChangePasswordSanitize() {
	req.CurrentPassword = strings.TrimSpace(req.CurrentPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}

type UpdateUserRequest struct {
	FirstName    *string `json:"first_name,omitempty" validate:"omitempty,min=3,max=50"`
	LastName     *string `json:"last_name,omitempty" validate:"omitempty,min=3,max=50"`
	Username     string  `json:"username" validate:"required,min=7,max=20"`
	Email        string  `json:"email" validate:"required,email,max=255"`
	MobileNumber string  `json:"mobile_number" validate:"required,min=5,max=255"`
}

func (req *UpdateUserRequest) UpdateUserSanitize() {
	if req.FirstName != nil {
		*req.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		*req.LastName = strings.TrimSpace(*req.LastName)
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
}

type SendPasswordResetEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (req *SendPasswordResetEmailRequest) SendPasswordResetEmailSanitize() {
	req.Email = strings.TrimSpace(req.Email)
}

type PasswordResetRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=7,max=60"`
}

func (req *PasswordResetRequest) PasswordResetSanitize() {
	req.Token = strings.TrimSpace(req.Token)
	req.Password = strings.TrimSpace(req.Password)
}

type ConfigureMFARequest struct {
	Email   string `json:"email" validate:"required,email,max=255"`
	Secret  string `json:"secret" validate:"required"`
	MFACode string `json:"mfa_code" validate:"required,len=6"`
}

func (req *ConfigureMFARequest) ConfigureMFASanitize() {
	req.Email = strings.TrimSpace(req.Email)
	req.Secret = strings.TrimSpace(req.Secret)
	req.MFACode = strings.TrimSpace(req.MFACode)
}

type VerifyMFARequest struct {
	Email   string `json:"email" validate:"required,email,max=255"`
	MFACode string `json:"mfa_code" validate:"required,len=6"`
}

func (req *VerifyMFARequest) VerifyMFASanitize() {
	req.Email = strings.TrimSpace(req.Email)
	req.MFACode = strings.TrimSpace(req.MFACode)
}
