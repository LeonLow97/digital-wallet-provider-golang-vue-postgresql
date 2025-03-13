package testdata

import (
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
)

// NewBalance returns a new instance of Balance with sample data.
func NewBalance() domain.Balance {
	return domain.Balance{
		ID:        1,
		Balance:   100.0,
		Currency:  "USD",
		UserID:    1,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

// NewBalanceWithID returns a new instance of Balance with the given ID and sample data.
func NewBalanceWithID(id int) domain.Balance {
	return domain.Balance{
		ID:        id,
		Balance:   150.0,
		Currency:  "EUR",
		UserID:    2,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

// NewBalances returns a slice of Balance with sample data.
func NewBalances() []domain.Balance {
	return []domain.Balance{
		{
			ID:        1,
			Balance:   100.0,
			Currency:  "USD",
			UserID:    1,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        2,
			Balance:   50.0,
			Currency:  "EUR",
			UserID:    2,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        3,
			Balance:   200.0,
			Currency:  "SGD",
			UserID:    1,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}
}

func NewWalletCurrencyAmount() []domain.WalletCurrencyAmount {
	return []domain.WalletCurrencyAmount{
		{
			WalletID:  1,
			Amount:    100.50,
			Currency:  "SGD",
			CreatedAt: "2024-06-16T12:00:00Z",
			UpdatedAt: "2024-06-16T12:00:00Z",
		},
		{
			WalletID:  2,
			Amount:    200.25,
			Currency:  "USD",
			CreatedAt: "2024-06-15T12:00:00Z",
			UpdatedAt: "2024-06-15T12:00:00Z",
		},
		{
			WalletID:  3,
			Amount:    50.00,
			Currency:  "EUR",
			CreatedAt: "2024-06-14T12:00:00Z",
			UpdatedAt: "2024-06-14T12:00:00Z",
		},
	}
}
