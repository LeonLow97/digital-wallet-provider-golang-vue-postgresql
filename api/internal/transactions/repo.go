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
	GetDB() *sqlx.DB

	GetUserCountByUserId(ctx context.Context, userId int) (int, error)
	GetUserIdByMobileNumber(ctx context.Context, mobileNumber string) (int, error)
	GetCountByUserIdAndBeneficiaryId(ctx context.Context, userId, beneficiaryId int) (int, error)

	GetCountByUserIdAndCurrency(tx *sql.Tx, ctx context.Context, userId int, currency string) (int, int, error)
	GetBalanceIdByUserIdAndPrimary(tx *sql.Tx, ctx context.Context, userId int) (int, string, error)
	GetBalanceAmountById(tx *sql.Tx, ctx context.Context, balanceId int) (float64, error)

	UpdateBalanceAmountById(tx *sql.Tx, ctx context.Context, balance float64, balanceId int) error
	InsertIntoTransactions(tx *sql.Tx, ctx context.Context, transaction *TransactionEntity) error 

	SQLTransactionMoneyTransfer(ctx context.Context, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus string, amountTransferred, amountReceived float64) error
	GetTransactionsCountByUserId(ctx context.Context, userId int) (int, error)
	GetTransactionsByUserId(ctx context.Context, userId, pageSize, offset int) (*Transactions, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetDB() *sqlx.DB {
	return r.db
}

// retrieve the user count by the specified userId
func (r *repo) GetUserCountByUserId(ctx context.Context, userId int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE id = $1;`

	if err := r.db.QueryRowContext(ctx, query, userId).Scan(&count); err != nil {
		return 0, utils.InternalServerError{Message: err.Error()}
	}
	return count, nil
}

// retrieve userId by the specified mobile number
func (r *repo) GetUserIdByMobileNumber(ctx context.Context, mobileNumber string) (int, error) {
	var userId int
	query := `SELECT id FROM users where mobile_number = $1;`

	if err := r.db.QueryRowContext(ctx, query, mobileNumber).Scan(&userId); err != nil {
		return 0, utils.InternalServerError{Message: err.Error()}
	}

	return userId, nil
}

func (r *repo) GetCountByUserIdAndBeneficiaryId(ctx context.Context, userId, beneficiaryId int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_beneficiary
		WHERE user_id = $1 AND beneficiary_id = $2;	
	`

	if err := r.db.QueryRowContext(ctx, query, userId, beneficiaryId).Scan(&count); err != nil {
		return 0, utils.InternalServerError{Message: err.Error()}
	}

	return count, nil
}

func (r *repo) GetCountByUserIdAndCurrency(tx *sql.Tx, ctx context.Context, userId int, currency string) (int, int, error) {
	var count, id int
	query := `SELECT COUNT(*), id FROM user_balance
		WHERE user_id = $1 AND currency = $2
		GROUP BY (id);
	`

	if err := tx.QueryRowContext(ctx, query, userId, currency).Scan(&count, &id); err != nil {
		return 0, 0, utils.InternalServerError{Message: err.Error()}
	}

	return count, id, nil
}

// retrieve primary balance by specified userId
func (r *repo) GetBalanceIdByUserIdAndPrimary(tx *sql.Tx, ctx context.Context, userId int) (int, string, error) {
	var id int
	var currency string
	query := `SELECT id, currency FROM user_balance
		WHERE user_id = $1 AND is_primary = 1 returning id;
	`

	if err := tx.QueryRowContext(ctx, query, userId).Scan(&id, &currency); err != nil {
		return 0, "", utils.InternalServerError{Message: err.Error()}
	}

	return id, currency, nil
}

// retrieve user balance by balanceId
func (r *repo) GetBalanceAmountById(tx *sql.Tx, ctx context.Context, balanceId int) (float64, error) {
	var balance float64
	query := `SELECT balance FROM user_balance
		WHERE id = $1;
	`

	if err := tx.QueryRowContext(ctx, query, balanceId).Scan(&balance); err != nil {
		return 0, utils.InternalServerError{Message: err.Error()}
	}

	return balance, nil
}

// update the balance of the specified balance id
func (r *repo) UpdateBalanceAmountById(tx *sql.Tx, ctx context.Context, balance float64, balanceId int) error {
	query := `UPDATE user_balance SET balance = $1
		WHERE id = $2;
	`

	_, err := tx.ExecContext(ctx, query, balance, balanceId)
	if err != nil {
		return utils.InternalServerError{Message: err.Error()}
	}

	return nil
}

// insert a transaction into transactions table
func (r *repo) InsertIntoTransactions(tx *sql.Tx, ctx context.Context, transaction *TransactionEntity) error {
	query := `INSERT INTO transactions (
		user_id, sender_id, beneficiary_id, transferred_amount, transferred_amount_currency,
		received_amount, received_amount_currency, status, transferred_date, received_date
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	_, err := tx.ExecContext(ctx, query,
		&transaction.UserId,
		&transaction.SenderId,
		&transaction.BeneficiaryId,
		&transaction.TransferredAmount,
		&transaction.TransferredAmountCurrency,
		&transaction.ReceivedAmount,
		&transaction.ReceivedAmountCurrency,
		&transaction.Status,
		time.Now(), time.Now(),
	)
	if err != nil {
		return utils.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (r *repo) SQLTransactionMoneyTransfer(ctx context.Context, senderName, beneficiaryName, amountTransferredCurrency, amountReceivedCurrency, confirmedStatus, receivedStatus string, amountTransferred, amountReceived float64) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return utils.InternalServerError{Message: "error in SQL Transaction for Money Transfer: " + err.Error()}
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
		return utils.InternalServerError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
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
		return utils.InternalServerError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
	}

	beneficiaryBalance := nullableBeneficiaryBalance.Float64
	finalBeneficiaryBalance := beneficiaryBalance + amountReceived

	if _, err := tx.ExecContext(ctx, queryUpdateUserBalance, finalUserBalance, senderName); err != nil {
		return utils.InternalServerError{Message: "error with ExecContext in updating sender balance: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryUpdateUserBalance, finalBeneficiaryBalance, beneficiaryName); err != nil {
		return utils.InternalServerError{Message: "error with ExecContext in updating beneficiary balance: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryInsertIntoTransactions, senderName, senderName, beneficiaryName, amountTransferred, amountTransferredCurrency, amountReceived, amountReceivedCurrency, confirmedStatus, time.Now(), time.Now()); err != nil {
		return utils.InternalServerError{Message: "error in ExecContext in inserting sender transaction: " + err.Error()}
	}

	if _, err := tx.ExecContext(ctx, queryInsertIntoTransactions, beneficiaryName, senderName, beneficiaryName, amountTransferred, amountTransferredCurrency, amountReceived, amountReceivedCurrency, receivedStatus, time.Now(), time.Now()); err != nil {
		return utils.InternalServerError{Message: "error in ExecContext in InsertIntoTransactions: " + err.Error()}
	}

	return nil
}

// Retrieves the number of transactions related to the specified userId
func (r *repo) GetTransactionsCountByUserId(ctx context.Context, userId int) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1;`

	var count int

	if err := r.db.QueryRowContext(ctx, query, userId).Scan(&count); err != nil {
		return 0, utils.InternalServerError{Message: err.Error()}
	}

	return count, nil
}

// paginated data - returns a list of transactions that is associated with the userId
func (r *repo) GetTransactionsByUserId(ctx context.Context, userId, pageSize, offset int) (*Transactions, error) {
	query := `
		SELECT
			COALESCE(s.first_name, '') AS sender_first_name,
			COALESCE(s.last_name, '') AS sender_last_name,
			s.username AS sender_username,
			COALESCE(b.first_name, '') AS beneficiary_first_name,
			COALESCE(b.last_name, '') AS beneficiary_last_name,
			b.username AS beneficiary_username,
			t.transferred_amount,
			t.transferred_amount_currency,
			t.received_amount,
			t.received_amount_currency,
			t.status,
			t.transferred_date,
			t.received_date
		FROM transactions t
		JOIN users s ON s.id = t.sender_id
		JOIN users b ON b.id = t.beneficiary_id
		WHERE t.user_id = $1
		ORDER BY t.transferred_date DESC
		LIMIT $2 OFFSET $3;
	`

	var transactions Transactions

	rows, err := r.db.QueryContext(ctx, query, userId, pageSize, offset)
	if err != nil {
		return nil, utils.InternalServerError{Message: err.Error()}
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.SenderFirstName,
			&transaction.SenderLastName,
			&transaction.SenderUsername,
			&transaction.BeneficiaryFirstName,
			&transaction.BeneficiaryLastName,
			&transaction.BeneficiaryUsername,
			&transaction.TransferredAmount,
			&transaction.TransferredAmountCurrency,
			&transaction.ReceivedAmount,
			&transaction.ReceivedAmountCurrency,
			&transaction.Status,
			&transaction.TransferredDate,
			&transaction.ReceivedDate,
		); err != nil {
			return nil, utils.InternalServerError{Message: err.Error()}
		}

		transactions.Transactions = append(transactions.Transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.InternalServerError{Message: err.Error()}
	}

	return &transactions, nil
}
