package dto

import "strings"

type CreateTransactionRequest struct {
	SenderWalletID      int     `json:"sender_wallet_id" validate:"required,min=1"`
	MobileNumber        string  `json:"mobile_number" validate:"required,min=5,max=255"`
	SourceCurrency      string  `json:"source_currency" validate:"required,min=3,max=3"`
	DestinationCurrency string  `json:"destination_currency" validate:"required,min=3,max=3"`
	Amount              float64 `json:"amount"`
}

func (req *CreateTransactionRequest) Sanitize() {
	req.MobileNumber = strings.TrimSpace(req.MobileNumber)
	req.SourceCurrency = strings.TrimSpace(req.SourceCurrency)
	req.DestinationCurrency = strings.TrimSpace(req.DestinationCurrency)
}
