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

func (r *walletRepository) GetWalletByWalletID(ctx context.Context, userID, walletID int) (*domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT
			w.id AS id,
			w.user_id AS user_id,
			wt.type AS wallet_type,
			wt.id AS wallet_type_id,
			w.created_at
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1 AND w.id = $2;
	`

	var wallet domain.Wallet
	if err := r.db.GetContext(ctx, &wallet, query, userID, walletID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNoWalletFound
		}
		return nil, err
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

	if len(wallets) == 0 {
		return nil, exception.ErrNoWalletsFound
	}

	return wallets, nil
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

	if len(walletBalances) == 0 {
		return nil, exception.ErrWalletBalancesNotFound
	}

	return walletBalances, nil
}

func (r *walletRepository) GetWalletBalancesByUserIDAndWalletID(ctx context.Context, tx *sqlx.Tx, userID, walletID int) ([]domain.WalletCurrencyAmount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT wallet_id, amount, currency, created_at, updated_at
		FROM wallet_balances
		WHERE user_id = $1 AND wallet_id = $2;
	`

	var walletBalances []domain.WalletCurrencyAmount
	var err error
	if tx == nil {
		err = r.db.SelectContext(ctx, &walletBalances, query, userID, walletID)
	} else {
		err = tx.SelectContext(ctx, &walletBalances, query, userID, walletID)
	}

	if err != nil {
		return nil, err
	}

	if len(walletBalances) == 0 {
		return nil, exception.ErrWalletBalancesNotFound
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

	if len(walletTypes) == 0 {
		return nil, exception.ErrWalletTypesNotFound
	}

	return &walletTypes, nil
}

func (r *walletRepository) CheckWalletExistsByWalletTypeID(ctx context.Context, userID, walletTypeID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM wallets
			WHERE 
				user_id = $1 AND
				wallet_type_id = $2
		)
	`

	var walletExists bool
	if err := r.db.QueryRowContext(ctx, query, userID, walletTypeID).Scan(&walletExists); err != nil {
		return false, err
	}

	return walletExists, nil
}

func (r *walletRepository) CheckWalletTypeExists(ctx context.Context, walletTypeID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM wallet_types
			WHERE id = $1
		)
	`

	var walletTypeExists bool
	if err := r.db.QueryRowContext(ctx, query, walletTypeID).Scan(&walletTypeExists); err != nil {
		return false, err
	}

	return walletTypeExists, nil
}

func (r *walletRepository) CreateWallet(ctx context.Context, tx *sqlx.Tx, wallet *domain.Wallet) (int, error) {
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

func (r *walletRepository) InsertWalletCurrencyAmount(ctx context.Context, tx *sqlx.Tx, walletID, userID int, currencyAmount []domain.WalletCurrencyAmount) error {
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

func (r *walletRepository) TopUpWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	placeholders := make([]string, 0)
	args := make([]interface{}, 0)

	for currency, amount := range finalWalletBalancesMap {
		placeholder := "( ?, ?, ?, ?, ?, ? )"
		placeholders = append(placeholders, placeholder)
		args = append(args, amount, currency, walletID, userID, time.Now(), time.Now())
	}

	// Syntax for UPSERT in Postgres
	// https://stackoverflow.com/questions/36359440/postgresql-insert-on-conflict-update-upsert-use-all-excluded-values
	// `ON CONFLICT DO UPDATE` all rows will be locked when the action is taken
	query := fmt.Sprintf(`
		INSERT INTO wallet_balances (amount, currency, wallet_id, user_id, created_at, updated_at)
		VALUES 
			%s
		ON CONFLICT (currency, wallet_id, user_id)
		DO UPDATE SET
			amount = EXCLUDED.amount,
			updated_at = EXCLUDED.updated_at;
	`, strings.Join(placeholders, ", "))

	if _, err := tx.ExecContext(ctx, r.db.Rebind(query), args...); err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) CashOutWalletBalances(ctx context.Context, tx *sqlx.Tx, userID, walletID int, finalWalletBalancesMap map[string]float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for currency, amount := range finalWalletBalancesMap {
		querySelect := `
			SELECT * FROM wallet_balances
			WHERE
				user_id = ? AND 
				wallet_id = ? AND 
				currency = ?
			FOR UPDATE;
		`

		query := `
			UPDATE wallet_balances
			SET 
				amount = ?,
				updated_at = ?
			WHERE
				user_id = ? AND
				wallet_id = ? AND
				currency = ?
		`

		// locking the row for wallet balance
		if _, err := tx.ExecContext(ctx, r.db.Rebind(querySelect), userID, walletID, currency); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, r.db.Rebind(query), amount, time.Now(), userID, walletID, currency); err != nil {
			return err
		}
	}

	return nil
}
