package transactions

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/LeonLow97/internal/utils"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetCountByUsername(ctx context.Context, username string) (int, error)
	GetCountByUsernameAndBeneficiaryNumber(ctx context.Context, username, beneficiaryNumber string) (int, error)
	GetCountByBeneficiaryNameAndBeneficiaryNumber(ctx context.Context, beneficiaryName, beneficiaryNumber string) (int, error)
	GetUsernameByBeneficiaryNumber(ctx context.Context, beneficiaryNumber string) (string, error)
	GetCurrencyByBeneficiaryMobileNumber(ctx context.Context, beneficiaryNumber string) (string, error)
	GetUserBalanceByUsername(ctx context.Context, username string) (float64, error)
	// UpdateUserBalanceByUsername(ctx context.Context, finalAmount float64, username string) error
	// InsertIntoTransactions(ctx context.Context, username, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, status string, amountTransferred, amountReceived float64) error
	SQLTransactionMoneyTransfer(ctx context.Context, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus string, amountTransferred, amountReceived float64) error
	GetByUserId(ctx context.Context, username string) (*Transactions, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetCountByUsername(ctx context.Context, username string) (int, error) {
	var count int
	query := `SELECT COUNT(username) FROM users WHERE username = $1;`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, username).Scan(&count); err != nil {
		return 0, utils.RepositoryError{Message: "error with QueryRowContext in GetCountByUsername: " + err.Error()}
	}
	return count, nil
}

func (r *repo) GetCountByUsernameAndBeneficiaryNumber(ctx context.Context, username, beneficiaryNumber string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_beneficiary
				WHERE user_id = (SELECT id FROM users WHERE username = $1)
				AND beneficiary_id = (SELECT beneficiary_id FROM beneficiaries WHERE mobile_number = $2);`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, username, beneficiaryNumber).Scan(&count); err != nil {
		return 0, utils.RepositoryError{Message: "error with QueryRowContext in GetCountByUsernameAndBeneficiaryNumber: " + err.Error()}
	}
	return count, nil
}

func (r *repo) GetCountByBeneficiaryNameAndBeneficiaryNumber(ctx context.Context, beneficiaryName, beneficiaryNumber string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE mobile_number = $1 AND username = $2;`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, beneficiaryNumber, beneficiaryName).Scan(&count); err != nil {
		return 0, utils.RepositoryError{Message: "error with QueryRowContext in GetCountByMobileNumber: " + err.Error()}
	}
	return count, nil
}

func (r *repo) GetUsernameByBeneficiaryNumber(ctx context.Context, beneficiaryNumber string) (string, error) {
	var username sql.NullString
	query := `SELECT username FROM users WHERE mobile_number = $1;`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, beneficiaryNumber).Scan(&username); err != nil {
		return "", utils.RepositoryError{Message: "error with QueryRowContext in GetUsernameByBeneficiaryNumber: " + err.Error()}
	}
	return username.String, nil
}

func (r *repo) GetCurrencyByBeneficiaryMobileNumber(ctx context.Context, beneficiaryNumber string) (string, error) {
	var currency sql.NullString
	query := `SELECT currency FROM beneficiaries WHERE mobile_number = $1;`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, beneficiaryNumber).Scan(&currency); err != nil {
		return "", utils.RepositoryError{Message: "error with QueryRowContext in GetCurrencyByBeneficiaryMobileNumber: " + err.Error()}
	}
	return currency.String, nil
}

func (r *repo) GetUserBalanceByUsername(ctx context.Context, username string) (float64, error) {
	var userBalance sql.NullFloat64
	query := `SELECT balance FROM user_balance
			WHERE user_id = (SELECT id FROM users WHERE username = $1);`

	query = r.db.Rebind(query)

	if err := r.db.QueryRowContext(ctx, query, username).Scan(&userBalance); err != nil {
		return 0.0, utils.RepositoryError{Message: "error with QueryRowContext in GetUserBalanceByUsername: " + err.Error()}
	}
	return userBalance.Float64, nil
}

func (r *repo) SQLTransactionMoneyTransfer(ctx context.Context, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus string, amountTransferred, amountReceived float64) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return utils.RepositoryError{Message: "error in SQL Transaction for Money Transfer: " + err.Error()}
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("error committing transaction: %s", err)
		}
	}()

	var nullableUserBalance, nullableBeneficiaryBalance sql.NullFloat64

	queryGetUserBalanceByUsername := `SELECT balance FROM user_balance
			WHERE user_id = (SELECT id FROM users WHERE username = $1);`
	queryUpdateUserBalance := `UPDATE user_balance
						SET balance = $1
						WHERE user_id = (SELECT id FROM users WHERE username = $2);`

	queryInsertIntoTransactions := `INSERT INTO transactions (user_id, sender_id, beneficiary_id, amount_transferred, amount_transferred_currency, amount_received, amount_received_currency, status, date_transferred, date_received)
				VALUES (
					(SELECT id FROM users WHERE username = $1),
					(SELECT id FROM users WHERE username = $2),
					(SELECT beneficiary_id FROM beneficiaries WHERE beneficiary_name = $3),
					$4, $5, $6, $7, $8, $9, $10
				);`

	queryGetUserBalanceByUsername = r.db.Rebind(queryGetUserBalanceByUsername)
	queryUpdateUserBalance = r.db.Rebind(queryUpdateUserBalance)
	queryInsertIntoTransactions = r.db.Rebind(queryInsertIntoTransactions)

	// Check sender balance and deduct from sender balance
	if err := tx.QueryRowContext(ctx, queryGetUserBalanceByUsername, senderName).Scan(&nullableUserBalance); err != nil {
		return utils.RepositoryError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
	}

	senderBalance := nullableUserBalance.Float64
	if senderBalance == 0.0 {
		return utils.ServiceError{Message: "Account has 0 balance. Please top up."}
	}
	if amountTransferred > senderBalance {
		return utils.ServiceError{Message: "User does not have sufficient funds to make the transfer. Please top up."}
	}

	finalUserBalance := senderBalance - amountTransferred

	// Add to beneficiary balance
	if err := tx.QueryRowContext(ctx, queryGetUserBalanceByUsername, beneficiaryName).Scan(&nullableBeneficiaryBalance); err != nil {
		return utils.RepositoryError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
	}

	beneficiaryBalance := nullableBeneficiaryBalance.Float64
	finalBeneficiaryBalance := beneficiaryBalance + amountReceived

	if _, err := tx.ExecContext(ctx, queryUpdateUserBalance, finalUserBalance, senderName); err != nil {
		return utils.RepositoryError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryUpdateUserBalance, finalBeneficiaryBalance, beneficiaryName); err != nil {
		return utils.RepositoryError{Message: "error with ExecContext in updating beneficiary balance: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryInsertIntoTransactions, senderName, senderName, beneficiaryName, amountTransferred, amountTransferredCurrency, amountReceived, amountReceivedCurrency, confirmedStatus, time.Now(), time.Now()); err != nil {
		return utils.RepositoryError{Message: "error in ExecContext in inserting sender transaction: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryInsertIntoTransactions, beneficiaryName, senderName, beneficiaryName, amountTransferred, amountTransferredCurrency, amountReceived, amountReceivedCurrency, receivedStatus, time.Now(), time.Now()); err != nil {
		return utils.RepositoryError{Message: "error in ExecContext in InsertIntoTransactions: " + err.Error()}
	}

	return nil
}

func (r *repo) GetByUserId(ctx context.Context, username string) (*Transactions, error) {
	query := `SELECT u.username, b.beneficiary_name, t.amount_transferred, t.amount_transferred_currency, t.amount_received, 
			t.amount_received_currency, t.status, t.date_transferred, t.date_received
			FROM transactions t
			LEFT JOIN users u ON u.id = t.sender_id
			LEFT JOIN beneficiaries b ON b.beneficiary_id = t.beneficiary_id
			WHERE t.user_id = (SELECT id FROM users WHERE username = $1);`

	var SenderName, BeneficiaryName, AmountTransferredCurrency, AmountReceivedCurrency, Status sql.NullString
	var AmountTransferred, AmountReceived sql.NullFloat64
	var DateTransferred, DateReceived sql.NullTime

	query = r.db.Rebind(query)

	rows, err := r.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, utils.RepositoryError{Message: "error with QueryContext: " + err.Error()}
	}

	defer rows.Close()

	var transactions Transactions

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&SenderName, &BeneficiaryName, &AmountTransferred, &AmountTransferredCurrency, &AmountReceived, &AmountReceivedCurrency, &Status, &DateTransferred, &DateReceived); err != nil {
			return nil, utils.RepositoryError{Message: "error scanning results: " + err.Error()}
		}

		transaction = Transaction{
			SenderName:                SenderName.String,
			BeneficiaryName:           BeneficiaryName.String,
			AmountTransferred:         AmountTransferred.Float64,
			AmountTransferredCurrency: AmountTransferredCurrency.String,
			AmountReceived:            AmountReceived.Float64,
			AmountReceivedCurrency:    AmountReceivedCurrency.String,
			Status:                    Status.String,
			DateTransferred:           DateTransferred.Time,
			DateReceived:              DateReceived.Time,
		}

		transactions.Transactions = append(transactions.Transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.RepositoryError{Message: "error iterating over result rows: " + err.Error()}
	}

	return &transactions, nil
}
