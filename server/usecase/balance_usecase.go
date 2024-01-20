package usecase

import "github.com/LeonLow97/go-clean-architecture/domain"

type balanceUsecase struct {
	balanceRepository domain.BalanceRepository
}

func NewBalanceUsecase(balanceRepository domain.BalanceRepository) domain.BalanceUsecase {
	return &balanceUsecase{
		balanceRepository: balanceRepository,
	}
}
