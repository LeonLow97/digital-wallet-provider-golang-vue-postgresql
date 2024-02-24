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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
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
