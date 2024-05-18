package usecase

import (
	"context"
	"log"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
)

type balanceUsecase struct {
	balanceRepository domain.BalanceRepository
}

func NewBalanceUsecase(balanceRepository domain.BalanceRepository) domain.BalanceUsecase {
	return &balanceUsecase{
		balanceRepository: balanceRepository,
	}
}

func (uc *balanceUsecase) GetBalances(ctx context.Context, userID int) (*dto.GetBalancesResponse, error) {
	balances, err := uc.balanceRepository.GetBalances(ctx, userID)
	if err != nil {
		log.Printf("failed to get balances for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	var resp dto.GetBalancesResponse
	for _, b := range *balances {
		balance := dto.GetBalanceResponse{
			ID:        b.ID,
			Balance:   b.Balance,
			Currency:  b.Currency,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
		resp.Balances = append(resp.Balances, balance)
	}

	return &resp, nil
}

func (uc *balanceUsecase) Deposit(ctx context.Context, req dto.DepositRequest) (*dto.GetBalanceResponse, error) {
	// In a real-world scenario, connect via Go HTTP client to the user's credit card API
	// to retrieve the deposited amount. For the purpose of this project, we assume
	// a successful retrieval, and req.Balance represents the received amount.

	balance, err := uc.balanceRepository.GetBalanceByUserID(ctx, req.UserID, req.Currency)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return nil, err
	}

	if balance != nil {
		// user already has this balance, update the balance
		balance.Balance = balance.Balance + req.Balance

		if err := uc.balanceRepository.UpdateBalance(ctx, balance); err != nil {
			return nil, err
		}

		return &dto.GetBalanceResponse{
			Balance:  balance.Balance,
			Currency: balance.Currency,
		}, nil
	}

	if balance == nil {
		// user does not have this balance, insert the balance
		b := domain.Balance{
			Balance:  req.Balance,
			Currency: req.Currency,
			UserID:   req.UserID,
		}

		if err := uc.balanceRepository.CreateBalance(ctx, &b); err != nil {
			return nil, err
		}
		return &dto.GetBalanceResponse{
			Balance:  b.Balance,
			Currency: b.Currency,
		}, nil
	}

	return nil, exception.ErrInternalServerError
}

func (uc *balanceUsecase) Withdraw(ctx context.Context, req dto.WithdrawRequest) (*dto.GetBalanceResponse, error) {
	// In a real-world scenario:
	// Connect to the customer's credit card API to initiate a withdrawal.
	// Once the withdrawal is successful and the credit card is updated,
	// receive a success message from the credit card API. Subsequently,
	// update the user's balance via Apache Kafka to mitigate potential failures.

	balance, err := uc.balanceRepository.GetBalanceByUserID(ctx, req.UserID, req.Currency)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return nil, err
	}

	updatedBalance := balance.Balance - req.Balance
	if updatedBalance < 0 {
		return nil, exception.ErrInsufficientFunds
	}

	balance.Balance = updatedBalance
	err = uc.balanceRepository.UpdateBalance(ctx, balance)
	if err != nil {
		return nil, err
	}

	resp := dto.GetBalanceResponse{
		Balance:  updatedBalance,
		Currency: balance.Currency,
	}
	return &resp, nil
}
