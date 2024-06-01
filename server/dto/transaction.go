package dto

import "strings"

type CreateTransactionRequest struct {
	SenderWalletID               int     `json:"sender_wallet_id" validate:"required,min=1"`
	SourceCurrency               string  `json:"source_currency" validate:"required,min=3,max=3"`
	SourceAmount                 float64 `json:"source_amount" validate:"required,gt=0"`
	BeneficiaryMobileCountryCode string  `json:"beneficiary_mobile_country_code" validate:"required,min=1,max=5"`
	BeneficiaryMobileNumber      string  `json:"beneficiary_mobile_number" validate:"required,min=5,max=255"`
}

func (req *CreateTransactionRequest) Sanitize() {
	req.SourceCurrency = strings.TrimSpace(req.SourceCurrency)
	req.BeneficiaryMobileNumber = strings.TrimSpace(req.BeneficiaryMobileNumber)
	req.BeneficiaryMobileCountryCode = strings.TrimSpace(req.BeneficiaryMobileCountryCode)
}
