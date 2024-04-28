package repository

import (
	"context"
	"database/sql"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileNumber string) error {
	query := `
		SELECT 1
		FROM user_beneficiary ub
		JOIN users u
			ON u.id = ub.beneficiary_id
		WHERE user_id = $1 AND u.mobile_number = $2;
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, mobileNumber).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return exception.ErrUserNotLinkedToBeneficiary
		}
		return err
	}

	return nil
}

func (r *transactionRepository) InsertTransaction(ctx context.Context, userID int, senderID int, beneficiaryID int, transaction domain.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, sender_id, beneficiary_id, source_of_transfer, sent_amount,
			source_currency, received_amount, received_currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`

	_, err := r.db.ExecContext(ctx, query,
		userID,
		senderID,
		beneficiaryID,
		transaction.SourceOfTransfer,
		transaction.SentAmount,
		transaction.SourceCurrency,
		transaction.ReceivedAmount,
		transaction.ReceivedCurrency,
		transaction.Status,
	)

	return err
}

func (r *transactionRepository) GetTransactions(ctx context.Context, userID int) (*[]domain.Transaction, error) {
	query := `
		SELECT
			sender.username 			AS sender_username,
			sender.mobile_number 		AS sender_mobile_number,
			beneficiary.username 		AS beneficiary_username,
			beneficiary.mobile_number 	AS beneficiary_mobile_number,
			t.sent_amount,
			t.source_currency,
			t.received_amount,
			t.received_currency,
			t.source_of_transfer,
			t.status,
			t.created_at
		FROM transactions t
		JOIN users AS sender
			ON t.sender_id = sender.id
		JOIN users AS beneficiary
			ON t.beneficiary_id = beneficiary.id
		WHERE t.user_id = $1
		ORDER BY created_at DESC;
	`

	var transactions []domain.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNoTransactionsFound
		}
	}

	return &transactions, nil
}
