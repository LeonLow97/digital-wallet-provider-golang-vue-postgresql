package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type TransactionRepository struct {
	mock.Mock
}

func (m *TransactionRepository) CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx context.Context, userID int, mobileCountryCode, mobileNumber string) (int, bool, bool, error) {
	args := m.Called(ctx, userID, mobileCountryCode, mobileNumber)

	return args.Int(0), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *TransactionRepository) CheckValidityOfSenderIDAndWalletID(ctx context.Context, userID, walletID int) (bool, string, error) {
	args := m.Called(ctx, userID, walletID)

	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *TransactionRepository) InsertTransaction(ctx context.Context, tx *sqlx.Tx, userID int, transaction domain.Transaction) error {
	args := m.Called(ctx, tx, userID, transaction)
	return args.Error(0)
}

func (m *TransactionRepository) GetTotalTransactionsCount(ctx context.Context, userID int, paginator *pagination.Paginator) error {
	args := m.Called(ctx, userID, paginator)
	return args.Error(0)
}

func (m *TransactionRepository) GetTransactions(ctx context.Context, userID int, paginator *pagination.Paginator) (*[]domain.Transaction, error) {
	args := m.Called(ctx, userID, paginator)

	var transactions *[]domain.Transaction
	if v, ok := args.Get(0).(*[]domain.Transaction); ok {
		transactions = v
	}

	return transactions, args.Error(1)
}
