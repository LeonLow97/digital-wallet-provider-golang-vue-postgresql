package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
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

func (r *transactionRepository) CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileCountryCode, mobileNumber string) (int, bool, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT ub.beneficiary_id, u.active, u.is_mfa_configured
		FROM user_beneficiary ub
		JOIN users u
			ON u.id = ub.beneficiary_id
		WHERE
			ub.user_id = $1 AND
			u.mobile_country_code = $2 AND
			u.mobile_number = $3;
	`

	var beneficiaryID int
	var isBeneficiaryActive, isMFAConfigured bool
	if err := r.db.QueryRowContext(ctx, query, userID, mobileCountryCode, mobileNumber).Scan(&beneficiaryID, &isBeneficiaryActive, &isMFAConfigured); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, false, exception.ErrUserNotLinkedToBeneficiary
		}
		return 0, false, false, err
	}

	return beneficiaryID, isBeneficiaryActive, isMFAConfigured, nil
}

func (r *transactionRepository) CheckValidityOfSenderIDAndWalletID(ctx context.Context, userID, walletID int) (bool, string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM wallets
			WHERE
				user_id = $1 AND
				id = $2
		),
		wt.type
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1 AND w.id = $2;
	`

	var validSenderWallet bool
	var walletType string
	if err := r.db.QueryRowContext(ctx, query, userID, walletID).Scan(&validSenderWallet, &walletType); err != nil {
		return false, "", err
	}
	return validSenderWallet, walletType, nil
}

func (r *transactionRepository) InsertTransaction(ctx context.Context, tx *sqlx.Tx, userID int, transaction domain.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO transactions 
			(user_id, sender_id, beneficiary_id, source_of_transfer, source_amount,
			source_currency, destination_amount, destination_currency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	if _, err := tx.ExecContext(ctx, query,
		userID,
		transaction.SenderID,
		transaction.BeneficiaryID,
		transaction.SourceOfTransfer,
		transaction.SourceAmount,
		transaction.SourceCurrency,
		transaction.DestinationAmount,
		transaction.DestinationCurrency,
		transaction.Status,
		time.Now(),
	); err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) GetTransactions(ctx context.Context, userID int, paginator *pagination.Paginator) (*[]domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT
			sender.username 			AS sender_username,
			sender.mobile_number 		AS sender_mobile_number,
			beneficiary.username 		AS beneficiary_username,
			beneficiary.mobile_number 	AS beneficiary_mobile_number,
			t.source_amount,
			t.source_currency,
			t.destination_amount,
			t.destination_currency,
			t.source_of_transfer,
			t.status,
			t.created_at
		FROM transactions t
		JOIN users AS sender
			ON t.sender_id = sender.id
		JOIN users AS beneficiary
			ON t.beneficiary_id = beneficiary.id
		WHERE t.user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?;
	`

	args := []interface{}{userID, paginator.Limit(), paginator.Offset()}

	var transactions []domain.Transaction
	if err := r.db.SelectContext(ctx, &transactions, r.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, exception.ErrNoTransactionsFound
	}

	return &transactions, nil
}

func (r *transactionRepository) GetTotalTransactionsCount(ctx context.Context, userID int, paginator *pagination.Paginator) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT COUNT(1)
		FROM transactions t
		JOIN users AS sender
			ON t.sender_id = sender.id
		JOIN users AS beneficiary
			ON t.beneficiary_id = beneficiary.id
		WHERE t.user_id = $1
	`

	return r.db.QueryRowContext(ctx, query, userID).Scan(&paginator.TotalRecords)
}
