package testdata

import (
	"time"

	"github.com/LeonLow97/go-clean-architecture/domain"
)

func Wallet() *domain.Wallet {
	return &domain.Wallet{
		ID:           1,
		WalletType:   "Personal",
		WalletTypeID: 1,
		UserID:       1,
		CreatedAt:    time.Now().String(),
		CurrencyAmount: []domain.WalletCurrencyAmount{
			{
				WalletID:  1,
				Amount:    100,
				Currency:  "SGD",
				CreatedAt: time.Now().String(),
				UpdatedAt: time.Now().String(),
			},
		},
	}
}
