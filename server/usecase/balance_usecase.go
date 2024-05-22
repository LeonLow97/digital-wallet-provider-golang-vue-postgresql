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

func (uc *balanceUsecase) GetBalanceHistory(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceHistory, error) {
	balanceHistory, err := uc.balanceRepository.GetBalanceHistory(ctx, userID, balanceID)
	if err != nil {
		log.Printf("failed to get balance history for user id: %d, balance id: %d with error: %v\n", userID, balanceID, err)
		return nil, err
	}

	resp := &dto.GetBalanceHistory{
		BalanceHistory: *balanceHistory,
	}

	return resp, nil
}

func (uc *balanceUsecase) GetBalance(ctx context.Context, userID int, balanceID int) (*dto.GetBalanceResponse, error) {
	balance, err := uc.balanceRepository.GetBalance(ctx, userID, balanceID)
	if err != nil {
		log.Printf("failed to get balance for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	resp := dto.GetBalanceResponse{
		ID:        balance.ID,
		Balance:   balance.Balance,
		Currency:  balance.Currency,
		CreatedAt: balance.CreatedAt,
		UpdatedAt: balance.UpdatedAt,
	}

	return &resp, nil
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

func (uc *balanceUsecase) Deposit(ctx context.Context, req dto.DepositRequest) error {
	// In a real-world scenario, connect via Go HTTP client to the user's credit card API
	// to retrieve the deposited amount. For the purpose of this project, we assume
	// a successful retrieval, and req.Balance represents the received amount.

	currentBalance, err := uc.balanceRepository.GetBalance(ctx, req.UserID, 1)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return err
	}

	var updatedBalance *domain.Balance

	// Update the balance if it exists
	if currentBalance != nil {
		currentBalance.Balance += req.Balance

		if err := uc.balanceRepository.UpdateBalance(ctx, currentBalance); err != nil {
			return err
		}
		updatedBalance = currentBalance
	} else {
		// Create a new balance if it does not exist
		updatedBalance = &domain.Balance{
			Balance:  req.Balance,
			Currency: req.Currency,
			UserID:   req.UserID,
		}
		// user does not have this balance, insert the balance
		if err := uc.balanceRepository.CreateBalance(ctx, updatedBalance); err != nil {
			return err
		}
	}

	defer func() {
		err = uc.balanceRepository.CreateBalanceHistory(ctx, updatedBalance, req.Balance, "deposit")
	}()

	if err != nil {
		log.Printf("failed to create balance history with error: %v\n", err)
		return err
	}

	return nil
}

func (uc *balanceUsecase) Withdraw(ctx context.Context, req dto.WithdrawRequest) error {
	// In a real-world scenario:
	// Connect to the customer's credit card API to initiate a withdrawal.
	// Once the withdrawal is successful and the credit card is updated,
	// receive a success message from the credit card API. Subsequently,
	// update the user's balance via Apache Kafka to mitigate potential failures.

	currentBalance, err := uc.balanceRepository.GetBalance(ctx, req.UserID, 1)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return err
	}

	if req.Balance > currentBalance.Balance {
		return exception.ErrInsufficientFunds
	}

	if currentBalance != nil {
		currentBalance.Balance -= req.Balance
		if err := uc.balanceRepository.UpdateBalance(ctx, currentBalance); err != nil {
			return err
		}
	} else {
		return exception.ErrBalanceNotFound
	}

	defer func() {
		err = uc.balanceRepository.CreateBalanceHistory(ctx, currentBalance, req.Balance, "withdraw")
	}()

	if err != nil {
		log.Printf("failed to create balance history with error: %v\n", err)
		return err
	}

	return nil
}
