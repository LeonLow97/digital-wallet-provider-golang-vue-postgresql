package dto

import "strings"

type CreateBeneficiaryRequest struct {
	MobileCountryCode string `json:"mobile_country_code" validate:"required,min=1,max=5"`
	MobileNumber      string `json:"mobile_number" validate:"required,min=1,max=255"`
}

type UpdateBeneficiaryRequest struct {
	IsDeleted     int `json:"is_deleted" validate:"min=0,max=1"`
	BeneficiaryID int `json:"beneficiary_id" validate:"required,min=1"`
}

type GetBeneficiaryResponse struct {
	BeneficiaryID                int    `json:"beneficiaryID"`
	IsDeleted                    int    `json:"isDeleted"`
	BeneficiaryFirstName         string `json:"beneficiaryFirstName"`
	BeneficiaryLastName          string `json:"beneficiaryLastName"`
	BeneficiaryEmail             string `json:"beneficiaryEmail"`
	BeneficiaryUsername          string `json:"beneficiaryUsername"`
	IsActive                     int    `json:"active"`
	BeneficiaryMobileCountryCode string `json:"beneficiaryMobileCountryCode"`
	BeneficiaryMobileNumber      string `json:"beneficiaryMobileNumber"`
}

type GetBeneficiariesResponse struct {
	Beneficiaries []GetBeneficiaryResponse `json:"beneficiaries"`
}

func (req *CreateBeneficiaryRequest) CreateBeneficiarySanitize() {
	req.MobileCountryCode = strings.TrimSpace(req.MobileCountryCode)
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
}
