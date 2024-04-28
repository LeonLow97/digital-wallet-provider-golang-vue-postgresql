package dto

import "strings"

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
