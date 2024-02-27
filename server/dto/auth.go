package dto

import "strings"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=7,max=20"`
	Password string `json:"password" validate:"required,min=7,max=60"`
}

type LoginResponse struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	MobileNumber string `json:"mobileNumber"`
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
