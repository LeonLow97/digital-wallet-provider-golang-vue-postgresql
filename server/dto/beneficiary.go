package dto

import "strings"

type CreateBeneficiaryRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required,min=1,max=255"`
	UserID       int
}

type UpdateBeneficiaryRequest struct {
	IsDeleted     int `json:"is_deleted" validate:"min=0,max=1"`
	BeneficiaryID int `json:"beneficiary_id" validate:"required,min=1"`
	UserID        int
}

func (req *CreateBeneficiaryRequest) CreateBeneficiarySanitize() {
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
}
