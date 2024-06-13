package domain

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
)

type Transaction struct {
	SenderID                int     `db:"sender_id"`
	BeneficiaryID           int     `db:"beneficiary_id"`
	SenderUsername          string  `json:"sender_username" db:"sender_username"`
	SenderMobileNumber      string  `json:"sender_mobile_number" db:"sender_mobile_number"`
	BeneficiaryUsername     string  `json:"beneficiary_username" db:"beneficiary_username"`
	BeneficiaryMobileNumber string  `json:"beneficiary_mobile_number" db:"beneficiary_mobile_number"`
	SourceAmount            float64 `json:"source_amount" db:"source_amount"`
	SourceCurrency          string  `json:"source_currency" db:"source_currency"`
	DestinationAmount       float64 `json:"destination_amount" db:"destination_amount"`
	DestinationCurrency     string  `json:"destination_currency" db:"destination_currency"`
	SourceOfTransfer        string  `json:"source_of_transfer" db:"source_of_transfer"`
	Status                  string  `json:"status" db:"status"`
	CreatedAt               string  `json:"created_at" db:"created_at"`
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID int) error
	GetTransactions(ctx context.Context, userID int, paginator *pagination.Paginator) (*[]Transaction, error)
}

type TransactionRepository interface {
	CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileCountryCode, mobileNumber string) (int, bool, bool, error)
	CheckValidityOfSenderIDAndWalletID(ctx context.Context, userID, walletID int) (bool, string, error) // TODO: move to wallet repository?

	InsertTransaction(ctx context.Context, tx *sql.Tx, userID int, transaction Transaction) error

	GetTotalTransactionsCount(ctx context.Context, userID int, paginator *pagination.Paginator) error
	GetTransactions(ctx context.Context, userID int, paginator *pagination.Paginator) (*[]Transaction, error)
}
