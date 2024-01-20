package repository

import (
	"context"
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
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
