package mocks

import (
	"context"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/utils/pagination"
	"github.com/stretchr/testify/mock"
)

type TransactionUsecase struct {
	mock.Mock
}

type TransactionUsecaseReturnValues struct {
	CreateTransaction []interface{}
	GetTransactions   []interface{}
}

func (m *TransactionUsecase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID int) error {
	args := m.Called(ctx, req, userID)
	return args.Error(0)
}

func (m *TransactionUsecase) GetTransactions(ctx context.Context, userID int, paginator *pagination.Paginator) (*[]domain.Transaction, error) {
	args := m.Called(ctx, userID, paginator)

	var transactions *[]domain.Transaction
	if v, ok := args.Get(0).(*[]domain.Transaction); ok {
		transactions = v
	}

	return transactions, args.Error(1)
}
