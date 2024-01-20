package usecase

import "github.com/LeonLow97/go-clean-architecture/domain"

type transactionUsecase struct {
	transactionRepository domain.TransactionRepository
}

func NewTransactionUsecase(transactionRepository domain.TransactionRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		transactionRepository: transactionRepository,
	}
}
