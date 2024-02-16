package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
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
		&wallet.Type,
		&wallet.TypeID,
		&wallet.Balance,
		&wallet.Currency,
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
			w.balance,
			w.currency,
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
		&wallet.Type,
		&wallet.TypeID,
		&wallet.Balance,
		&wallet.Currency,
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
			wt.type AS type,
			wt.id AS type_id,
			w.balance,
			w.currency,
			w.created_at
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1;
	`

	var wallets []domain.Wallet
	err := r.db.SelectContext(ctx, &wallets, query, userID)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *walletRepository) GetWalletTypes(ctx context.Context) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT id, type FROM wallet_types;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// use a set to store the types
	walletTypes := make(map[string]int)

	for rows.Next() {
		var walletType string
		var walletID int

		if err := rows.Scan(
			&walletID,
			&walletType,
		); err != nil {
			return nil, err
		}
		walletTypes[walletType] = walletID
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return walletTypes, nil
}

func (r *walletRepository) CheckWalletExistsByUserID(ctx context.Context, userID int, walletType string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		SELECT COUNT(w.id)
		FROM wallets w
		JOIN wallet_types wt
			ON w.wallet_type_id = wt.id
		WHERE w.user_id = $1 AND wt.type = $2;
	`

	var count int
	if err := r.db.QueryRowContext(ctx, query, userID, walletType).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *walletRepository) GetBalanceByUserID(ctx context.Context, userID int) (*domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT balance, currency
		FROM balances
		WHERE user_id = $1;
	`

	var balance domain.Balance
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&balance.Balance,
		&balance.Currency,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrBalanceNotFound
		}
	}

	return &balance, nil
}

func (r *walletRepository) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO wallets (balance, currency, wallet_type_id, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.db.ExecContext(ctx, query, wallet.Balance, wallet.Currency, wallet.TypeID, wallet.UserID, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) UpdateWallet(ctx context.Context, wallet *domain.Wallet) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		UPDATE wallets
		SET balance = $1, updated_at = $2
		WHERE user_id = $3 AND wallet_type_id = $4;
	`

	_, err := r.db.ExecContext(ctx, query,
		wallet.Balance,
		time.Now(),
		wallet.UserID,
		wallet.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
