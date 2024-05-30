package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/jmoiron/sqlx"
)

type walletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) domain.WalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) GetWalletByWalletType(ctx context.Context, userID int, walletType string) (*domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT
			w.id AS id,
			w.user_id AS user_id,
			wt.type AS type,
			wt.id AS type_id,
			w.balance,
			w.currency,
			w.created_at
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1 AND wt.type = $2;
	`

	var wallet domain.Wallet
	err := r.db.QueryRowContext(ctx, query, userID, walletType).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletType,
		&wallet.WalletTypeID,
		&wallet.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNoWalletFound
		}
	}
	return &wallet, nil
}

func (r *walletRepository) GetWalletByWalletID(ctx context.Context, userID, walletID int) (*domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT
			w.id AS id,
			w.user_id AS user_id,
			wt.type AS type,
			wt.id AS type_id,
			w.created_at
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1 AND w.id = $2;
	`

	var wallet domain.Wallet
	err := r.db.QueryRowContext(ctx, query, userID, walletID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.WalletType,
		&wallet.WalletTypeID,
		&wallet.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNoWalletFound
		}
	}
	return &wallet, nil
}

func (r *walletRepository) GetWallets(ctx context.Context, userID int) ([]domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT
			w.id AS id,
			wt.type AS wallet_type,
			wt.id AS wallet_type_id,
			w.created_at
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1;
	`

	var wallets []domain.Wallet
	if err := r.db.SelectContext(ctx, &wallets, query, userID); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *walletRepository) GetWalletBalancesByUserID_TX(ctx context.Context, tx *sql.Tx, userID int) ([]domain.WalletCurrencyAmount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT wallet_id, amount, currency, created_at, updated_at
		FROM wallet_balances
		WHERE user_id = $1;
	`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var walletBalances []domain.WalletCurrencyAmount
	for rows.Next() {
		var balance domain.WalletCurrencyAmount
		if err := rows.Scan(&balance.WalletID, &balance.Amount, &balance.Currency, &balance.CreatedAt, &balance.UpdatedAt); err != nil {
			return nil, err
		}
		walletBalances = append(walletBalances, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return walletBalances, nil
}

func (r *walletRepository) GetWalletBalancesByUserID(ctx context.Context, userID int) ([]domain.WalletCurrencyAmount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT wallet_id, amount, currency, created_at, updated_at
		FROM wallet_balances
		WHERE user_id = $1;
	`

	var walletBalances []domain.WalletCurrencyAmount
	if err := r.db.SelectContext(ctx, &walletBalances, query, userID); err != nil {
		return nil, err
	}

	return walletBalances, nil
}

func (r *walletRepository) GetWalletBalancesByUserIDAndWalletID(ctx context.Context, userID, walletID int) ([]domain.WalletCurrencyAmount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT wallet_id, amount, currency, created_at, updated_at
		FROM wallet_balances
		WHERE user_id = $1 AND wallet_id = $2;
	`

	var walletBalances []domain.WalletCurrencyAmount
	if err := r.db.SelectContext(ctx, &walletBalances, query, userID, walletID); err != nil {
		return nil, err
	}

	return walletBalances, nil
}

func (r *walletRepository) GetWalletBalancesByUserIDAndWalletID_TX(ctx context.Context, tx *sql.Tx, userID, walletID int) ([]domain.WalletCurrencyAmount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT wallet_id, amount, currency, created_at, updated_at
		FROM wallet_balances
		WHERE user_id = $1 AND wallet_id = $2;
	`

	var walletBalances []domain.WalletCurrencyAmount
	rows, err := tx.QueryContext(ctx, query, userID, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var walletBalance domain.WalletCurrencyAmount
		if err := rows.Scan(&walletBalance.WalletID, &walletBalance.Amount, &walletBalance.Currency, &walletBalance.CreatedAt, &walletBalance.UpdatedAt); err != nil {
			return nil, err
		}
		walletBalances = append(walletBalances, walletBalance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return walletBalances, nil
}

func (r *walletRepository) GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT id, type FROM wallet_types;
	`

	var walletTypes []dto.GetWalletTypesResponse
	if err := r.db.SelectContext(ctx, &walletTypes, query); err != nil {
		return nil, err
	}

	return &walletTypes, nil
}

func (r *walletRepository) PerformWalletValidationByUserID(ctx context.Context, userID, walletTypeID int) (*dto.WalletValidation, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT 
		  (SELECT EXISTS (
			SELECT 1
			FROM wallets w
			JOIN wallet_types wt
			  ON w.wallet_type_id = wt.id
			WHERE w.user_id = $1 AND 
				wt.id = $2
		  )) AS wallet_exists,
		  (SELECT EXISTS (
			SELECT 1
			FROM wallet_types
			WHERE id = $3
		  )) AS is_valid_wallet_type
	`

	args := []interface{}{userID, walletTypeID, walletTypeID}

	var walletValidation dto.WalletValidation
	if err := r.db.GetContext(ctx, &walletValidation, query, args...); err != nil { // Pass the third argument "valid_wallet_type"
		return nil, err
	}
	return &walletValidation, nil
}

func (r *walletRepository) GetAllBalancesByUserID(ctx context.Context, tx *sql.Tx, userID int) ([]domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT amount, currency
		FROM wallet_balances
		WHERE user_id = $1;
	`

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []domain.Balance
	for rows.Next() {
		var balance domain.Balance
		if err := rows.Scan(&balance.Balance, &balance.Currency); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}

func (r *walletRepository) CreateWallet(ctx context.Context, tx *sql.Tx, wallet *domain.Wallet) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO wallets (wallet_type_id, user_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var walletID int
	if err := tx.QueryRowContext(ctx, query, wallet.WalletTypeID, wallet.UserID, time.Now()).Scan(&walletID); err != nil {
		return 0, err
	}

	return walletID, nil
}

func (r *walletRepository) InsertWalletCurrencyAmount(ctx context.Context, tx *sql.Tx, walletID, userID int, currencyAmount []domain.WalletCurrencyAmount) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO wallet_balances (amount, currency, wallet_id, user_id)
		VALUES ($1, $2, $3, $4);
	`

	for _, c := range currencyAmount {
		if _, err := tx.ExecContext(ctx, query, c.Amount, c.Currency, walletID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (r *walletRepository) TopUpWalletBalances(ctx context.Context, tx *sql.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	i := 1
	for currency, amount := range finalWalletBalancesMap {
		placeholder := fmt.Sprintf("( $%d, $%d, $%d, $%d, $%d, $%d )", i, i+1, i+2, i+3, i+4, i+5)
		placeholders = append(placeholders, placeholder)
		args = append(args, amount, currency, walletID, userID, time.Now(), time.Now())
		i += 6
	}

	// Syntax for UPSERT in Postgres
	// https://stackoverflow.com/questions/36359440/postgresql-insert-on-conflict-update-upsert-use-all-excluded-values
	query := fmt.Sprintf(`
		INSERT INTO wallet_balances (amount, currency, wallet_id, user_id, created_at, updated_at)
		VALUES 
			%s
		ON CONFLICT (currency, wallet_id, user_id)
		DO UPDATE SET
			amount = EXCLUDED.amount,
			updated_at = EXCLUDED.updated_at;
	`, strings.Join(placeholders, ", "))

	_, err := tx.ExecContext(ctx, query, args...)

	return err
}

func (r *walletRepository) CashOutWalletBalances(ctx context.Context, tx *sql.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for currency, amount := range finalWalletBalancesMap {
		query := `
			UPDATE wallet_balances
			SET 
				amount = $1,
				updated_at = $2
			WHERE
				user_id = $3 AND
				wallet_id = $4 AND
				currency = $5
		`

		_, err := tx.ExecContext(ctx, query, amount, time.Now(), userID, walletID, currency)
		if err != nil {
			return err
		}
	}

	return nil
}
