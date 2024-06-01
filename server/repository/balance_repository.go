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

type balanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) domain.BalanceRepository {
	return &balanceRepository{
		db: db,
	}
}

func (r *balanceRepository) GetBalanceHistory(ctx context.Context, userID, balanceID int) (*[]dto.BalanceHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, amount, currency, type, user_id, balance_id, created_at
		FROM balances_history
		WHERE user_id = $1 AND balance_id = $2;
	`

	var balanceHistory []dto.BalanceHistory
	if err := r.db.SelectContext(ctx, &balanceHistory, query, userID, balanceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrBalanceHistoryNotFound
		}
		return nil, err
	}
	return &balanceHistory, nil
}

func (r *balanceRepository) GetBalances(ctx context.Context, tx *sql.Tx, userID int) ([]domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, balance, currency, created_at, updated_at
		FROM balances
		WHERE user_id = $1;
	`

	var balances []domain.Balance
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var balance domain.Balance
		if err := rows.Scan(&balance.ID, &balance.Balance, &balance.Currency, &balance.CreatedAt, &balance.UpdatedAt); err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}

func (r *balanceRepository) GetBalance(ctx context.Context, userID int, currency string) (*domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, balance, currency, user_id, created_at, updated_at
		FROM balances
		WHERE user_id = $1 AND currency = $2;
	`

	var balance domain.Balance
	err := r.db.QueryRowContext(ctx, query, userID, currency).Scan(
		&balance.ID,
		&balance.Balance,
		&balance.Currency,
		&balance.UserID,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrBalanceNotFound
		}
		return nil, err
	}
	return &balance, nil
}

func (r *balanceRepository) GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT currency
		FROM balances
		WHERE user_id = $1;
	`

	var balanceCurrencies []dto.GetUserBalanceCurrenciesResponse
	if err := r.db.SelectContext(ctx, &balanceCurrencies, query, userID); err != nil {
		return nil, err
	}

	return &balanceCurrencies, nil
}

func (r *balanceRepository) GetBalanceTx(ctx context.Context, tx *sql.Tx, userID int, currency string) (*domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, balance, currency, user_id, created_at, updated_at
		FROM balances
		WHERE user_id = $1 AND currency = $2;
	`

	var balance domain.Balance
	err := tx.QueryRowContext(ctx, query, userID, currency).Scan(
		&balance.ID,
		&balance.Balance,
		&balance.Currency,
		&balance.UserID,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrBalanceNotFound
		}
		return nil, err
	}
	return &balance, nil
}

func (r *balanceRepository) GetBalanceById(ctx context.Context, userID int, balanceId int) (*domain.Balance, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
		SELECT id, balance, currency, user_id, created_at, updated_at
		FROM balances
		WHERE user_id = $1 AND id = $2;
	`

	var balance domain.Balance
	err := r.db.QueryRowContext(ctx, query, userID, balanceId).Scan(
		&balance.ID,
		&balance.Balance,
		&balance.Currency,
		&balance.UserID,
		&balance.CreatedAt,
		&balance.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrBalanceNotFound
		}
		return nil, err
	}
	return &balance, nil
}

func (r *balanceRepository) CreateBalance(ctx context.Context, tx *sql.Tx, balance *domain.Balance) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO balances (balance, currency, user_id)
		VALUES ($1, $2, $3);
	`

	_, err := tx.ExecContext(ctx, query,
		balance.Balance,
		balance.Currency,
		balance.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *balanceRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, balance *domain.Balance) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		UPDATE balances
		SET balance = $1, updated_at = NOW()
		WHERE user_id = $2 AND currency = $3;
	`

	_, err := tx.ExecContext(ctx, query,
		balance.Balance,
		balance.UserID,
		balance.Currency,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *balanceRepository) UpdateBalances(ctx context.Context, tx *sql.Tx, userID int, finalBalancesMap map[string]float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for currency, newBalance := range finalBalancesMap {
		query := `
			INSERT INTO balances (user_id, currency, balance)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, currency)
			DO UPDATE SET balance = EXCLUDED.balance;
		`

		_, err := tx.ExecContext(ctx, query, userID, currency, newBalance)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *balanceRepository) CreateBalanceHistory(ctx context.Context, tx *sql.Tx, balance *domain.Balance, depositedBalance float64, balanceType string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	query := `
		INSERT INTO balances_history (amount, currency, type, user_id, balance_id)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := tx.ExecContext(ctx, query,
		depositedBalance,
		balance.Currency,
		balanceType,
		balance.UserID,
		balance.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
