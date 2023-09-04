package transactions

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) GetDB() *sqlx.DB {
	return nil
}

func (m *mockRepo) GetUserCountByUserId(ctx context.Context, userId int) (int, error) {
	args := m.Called(ctx, userId)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) GetUserIdByMobileNumber(ctx context.Context, mobileNumber string) (int, error) {
	args := m.Called(ctx, mobileNumber)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) GetCountByUserIdAndBeneficiaryId(ctx context.Context, userId, beneficiaryId int) (int, error) {
	args := m.Called(ctx, userId, beneficiaryId)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) GetCountByUserIdAndCurrency(tx *sql.Tx, ctx context.Context, userId int, currency string) (int, int, error) {
	args := m.Called(tx, ctx, userId, currency)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *mockRepo) GetBalanceIdByUserIdAndPrimary(tx *sql.Tx, ctx context.Context, userId int) (int, string, error) {
	args := m.Called(tx, ctx, userId)
	return args.Int(0), args.String(1), args.Error(2)
}

func (m *mockRepo) GetBalanceAmountById(tx *sql.Tx, ctx context.Context, balanceId int) (float64, error) {
	args := m.Called(tx, ctx, balanceId)
	return float64(args.Int(0)), args.Error(1)
}

func (m *mockRepo) UpdateBalanceAmountById(tx *sql.Tx, ctx context.Context, balance float64, balanceId int) error {
	args := m.Called(tx, ctx, balance, balanceId)
	return args.Error(0)
}

func (m *mockRepo) InsertIntoTransactions(tx *sql.Tx, ctx context.Context, transaction *TransactionEntity) error {
	args := m.Called(tx, ctx, transaction)
	return args.Error(0)
}

func (m *mockRepo) GetTransactionsCountByUserId(ctx context.Context, userId int) (int, error) {
	args := m.Called(ctx, userId)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) GetTransactionsByUserId(ctx context.Context, userId, pageSize, offset int) (*Transactions, error) {
	args := m.Called(ctx, userId, pageSize, offset)

	// Check if the first return value (index 0) is a *Transactions type.
	var transactions *Transactions
	if t, ok := args.Get(0).(*Transactions); ok {
		transactions = t
	}

	return transactions, args.Error(1)
}
