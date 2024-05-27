package repository

import (
	"context"
	"database/sql"
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
			WHERE w.user_id = $1 AND wt.id = $2
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

func (r *walletRepository) GetAllBalancesByUserID(ctx context.Context, userID int) ([]domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT balance, currency
		FROM balances
		WHERE user_id = $1;
	`

	var balances []domain.Balance
	if err := r.db.SelectContext(ctx, &balances, query, userID); err != nil {
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

func (r *walletRepository) UpdateWallet(ctx context.Context, wallet *domain.Wallet) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		UPDATE wallets
		SET balance = $1, updated_at = NOW()
		WHERE user_id = $2 AND wallet_type_id = $3 AND currency = $4;
	`

	_, err := r.db.ExecContext(ctx, query,
		wallet.UserID,
		wallet.WalletTypeID,
	)
	if err != nil {
		return err
	}
	return nil
}
