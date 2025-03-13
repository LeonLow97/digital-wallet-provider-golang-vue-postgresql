package testdata

import (
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
)

var fixedTime = "2024-06-15T13:45:00Z"

// MockWallet returns a mock instance of *domain.Wallet with static data.
func MockWallet() *domain.Wallet {
	return &domain.Wallet{
		ID:           1,
		WalletType:   "Personal",
		WalletTypeID: 1,
		UserID:       1,
		CreatedAt:    fixedTime,
	}
}

// MockWalletCurrencyAmounts returns a slice of domain.WalletCurrencyAmount with static data.
func MockWalletCurrencyAmounts() []domain.WalletCurrencyAmount {
	return []domain.WalletCurrencyAmount{
		{
			WalletID:  1,
			Amount:    100.0,
			Currency:  "SGD",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			WalletID:  1,
			Amount:    50.0,
			Currency:  "USD",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}
}

// MockWalletWithCurrencyAmounts returns a mock instance of *domain.Wallet with CurrencyAmount initialized.
func MockWalletWithCurrencyAmounts() *domain.Wallet {
	return &domain.Wallet{
		ID:             1,
		WalletType:     "Personal",
		WalletTypeID:   1,
		UserID:         1,
		CreatedAt:      fixedTime,
		CurrencyAmount: MockWalletCurrencyAmounts(),
	}
}

// MockWallets returns a slice of *domain.Wallet pointers with count number of mock instances.
func MockWallets(count int) []domain.Wallet {
	wallets := make([]domain.Wallet, 0, count)
	for i := 0; i < count; i++ {
		wallet := domain.Wallet{
			ID:           i + 1,
			WalletType:   "Personal",
			WalletTypeID: 1,
			UserID:       1,
			CreatedAt:    fixedTime,
		}
		wallets = append(wallets, wallet)
	}
	return wallets
}

// MockWalletsCurrencyAmounts returns mock data for WalletCurrencyAmount with the given walletID.
func MockWalletCurrencyAmountsByWalletID(walletID int) []domain.WalletCurrencyAmount {
	return []domain.WalletCurrencyAmount{
		{
			WalletID:  walletID,
			Amount:    float64(walletID * 100),
			Currency:  "SGD",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
		{
			WalletID:  walletID,
			Amount:    float64(walletID * 50),
			Currency:  "USD",
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}
}

// MockGetWalletTypesResponses generates mock data for []*GetWalletTypesResponse.
func MockGetWalletTypesResponses() *[]dto.GetWalletTypesResponse {
	return &[]dto.GetWalletTypesResponse{
		{
			ID:         1,
			WalletType: "Personal",
		},
		{
			ID:         2,
			WalletType: "Business",
		},
	}
}
