package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Transaction struct {
	SenderUsername          string  `json:"sender_username" db:"sender_username"`
	SenderMobileNumber      string  `json:"sender_mobile_number" db:"sender_mobile_number"`
	BeneficiaryUsername     string  `json:"beneficiary_username" db:"beneficiary_username"`
	BeneficiaryMobileNumber string  `json:"beneficiary_mobile_number" db:"beneficiary_mobile_number"`
	SentAmount              float64 `json:"sent_amount" db:"sent_amount"`
	SourceCurrency          string  `json:"source_currency" db:"source_currency"`
	ReceivedAmount          float64 `json:"received_amount" db:"received_amount"`
	ReceivedCurrency        string  `json:"received_currency" db:"received_currency"`
	SourceOfTransfer        string  `json:"source_of_transfer" db:"source_of_transfer"`
	Status                  string  `json:"status" db:"status"`
	CreatedAt               string  `json:"created_at" db:"created_at"`
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID int) error
	GetTransactions(ctx context.Context, userID int) (*[]Transaction, error)
}

type TransactionRepository interface {
	CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileNumber string) error
	InsertTransaction(ctx context.Context, userID int, senderID int, beneficiaryID int, transaction Transaction) error
	GetTransactions(ctx context.Context, userID int) (*[]Transaction, error)
}
