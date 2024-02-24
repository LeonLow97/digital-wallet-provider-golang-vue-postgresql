package domain

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/dto"
)

type Transaction struct {
	SenderUsername          string  `db:"sender_username"`
	SenderMobileNumber      string  `db:"sender_mobile_number"`
	BeneficiaryUsername     string  `db:"beneficiary_username"`
	BeneficiaryMobileNumber string  `db:"beneficiary_mobile_number"`
	SentAmount              float64 `db:"sent_amount"`
	SourceCurrency          string  `db:"source_currency"`
	ReceivedAmount          float64 `db:"received_amount"`
	ReceivedCurrency        string  `db:"received_currency"`
	SourceOfTransfer        string  `db:"source_of_transfer"`
	Status                  string  `db:"status"`
	CreatedAt               string  `db:"created_at"`
}

type TransactionUsecase interface {
	CreateTransactionByWallet(ctx context.Context, req dto.CreateTransactionRequest, userID int) error
}

type TransactionRepository interface {
	CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileNumber string) error
	InsertTransaction(ctx context.Context, userID int, senderID int, beneficiaryID int, transaction Transaction) error
}
