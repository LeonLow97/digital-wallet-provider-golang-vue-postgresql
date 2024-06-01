package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/LeonLow97/go-clean-architecture/utils"
	"github.com/jmoiron/sqlx"
)

type balanceUsecase struct {
	dbConn            *sqlx.DB
	userRepository    domain.UserRepository
	balanceRepository domain.BalanceRepository
}

func NewBalanceUsecase(dbConn *sqlx.DB, userRepository domain.UserRepository, balanceRepository domain.BalanceRepository) domain.BalanceUsecase {
	return &balanceUsecase{
		dbConn:            dbConn,
		userRepository:    userRepository,
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
	balance, err := uc.balanceRepository.GetBalanceById(ctx, userID, balanceID)
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
	// Start SQL Transaction, need to lock balance in case use POSTMAN and frontend to update balance at the same time
	tx, err := uc.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to begin sql transaction with error: %v\n", err)
		return nil, err
	}

	// Defer rollback or commit the transaction based on the outcome
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
			log.Printf("failed to complete sql transaction with error: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("failed to commit sql transaction with error: %v\n", err)
			}
		}
	}()

	balances, err := uc.balanceRepository.GetBalances(ctx, tx, userID)
	if err != nil {
		log.Printf("failed to get balances for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	var resp dto.GetBalancesResponse
	for _, b := range balances {
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

func (uc *balanceUsecase) GetUserBalanceCurrencies(ctx context.Context, userID int) (*[]dto.GetUserBalanceCurrenciesResponse, error) {
	resp, err := uc.balanceRepository.GetUserBalanceCurrencies(ctx, userID)
	if err != nil {
		log.Printf("failed to get user balance currencies for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	return resp, nil
}

func (uc *balanceUsecase) Deposit(ctx context.Context, req dto.DepositRequest) error {
	// In a real-world scenario, connect via Go HTTP client to the user's credit card API
	// to retrieve the deposited amount. For the purpose of this project, we assume
	// a successful retrieval, and req.Balance represents the received amount.

	// Start SQL Transaction, need to lock balance in case use POSTMAN and frontend to update balance at the same time
	tx, err := uc.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to begin sql transaction with error: %v\n", err)
		return err
	}

	// Defer rollback or commit the transaction based on the outcome
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
			log.Printf("failed to complete sql transaction with error: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("failed to commit sql transaction with error: %v\n", err)
			}
		}
	}()

	// retrieve user mobile country code (assumption: to determine the currency the user can deposit)
	user, err := uc.userRepository.GetUserByID(ctx, req.UserID)
	if err != nil {
		log.Printf("failed to get user with error: %v\n", err)
		return err
	}
	if utils.CountryCodeToCurrencyMap[user.MobileCountryCode] != req.Currency {
		return exception.ErrDepositCurrencyNotAllowed
	}

	currentBalance, err := uc.balanceRepository.GetBalanceTx(ctx, tx, req.UserID, req.Currency)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return err
	}

	var updatedBalance *domain.Balance

	// Update the balance if it exists
	if currentBalance != nil {
		currentBalance.Balance += req.Balance

		if err := uc.balanceRepository.UpdateBalance(ctx, tx, currentBalance); err != nil {
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
		if err := uc.balanceRepository.CreateBalance(ctx, tx, updatedBalance); err != nil {
			return err
		}
	}

	defer func() {
		err = uc.balanceRepository.CreateBalanceHistory(ctx, tx, updatedBalance, req.Balance, "deposit")
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

	// Start SQL Transaction, need to lock balance in case use POSTMAN and frontend to update balance at the same time
	tx, err := uc.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to begin sql transaction with error: %v\n", err)
		return err
	}

	// Defer rollback or commit the transaction based on the outcome
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
			log.Printf("failed to complete sql transaction with error: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("failed to commit sql transaction with error: %v\n", err)
			}
		}
	}()

	// retrieve user mobile country code (assumption: to determine the currency the user can deposit)
	user, err := uc.userRepository.GetUserByID(ctx, req.UserID)
	if err != nil {
		log.Printf("failed to get user with error: %v\n", err)
		return err
	}
	if utils.CountryCodeToCurrencyMap[user.MobileCountryCode] != req.Currency {
		return exception.ErrWithdrawCurrencyNotAllowed
	}

	currentBalance, err := uc.balanceRepository.GetBalanceTx(ctx, tx, req.UserID, req.Currency)
	if err != nil {
		log.Printf("failed to get one balance for user id %d with error: %v\n", req.UserID, err)
		return err
	}

	if req.Balance > currentBalance.Balance {
		return exception.ErrInsufficientFunds
	}

	if currentBalance != nil {
		currentBalance.Balance -= req.Balance
		if err := uc.balanceRepository.UpdateBalance(ctx, tx, currentBalance); err != nil {
			return err
		}
	} else {
		return exception.ErrBalanceNotFound
	}

	defer func() {
		err = uc.balanceRepository.CreateBalanceHistory(ctx, tx, currentBalance, req.Balance, "withdraw")
	}()

	if err != nil {
		log.Printf("failed to create balance history with error: %v\n", err)
		return err
	}

	return nil
}

func (uc *balanceUsecase) CurrencyExchange(ctx context.Context, userID int, req dto.CurrencyExchangeRequest) error {
	// Start SQL Transaction, need to lock balance in case use POSTMAN and frontend to update balance at the same time
	tx, err := uc.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to begin sql transaction with error: %v\n", err)
		return err
	}

	// Defer rollback or commit the transaction based on the outcome
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
			log.Printf("failed to complete sql transaction with error: %v\n", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Printf("failed to commit sql transaction with error: %v\n", err)
			}
		}
	}()

	user, err := uc.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("failed to get user with error: %v\n", err)
		return err
	}

	fromCurrency := utils.CountryCodeToCurrencyMap[user.MobileCountryCode]

	// TODO: check if toCurrencies is allowed with a list of currencies stored in the database, for now use a map
	if _, found := utils.ToCurrencies[req.ToCurrency]; !found {
		log.Printf("toCurrency %s is not under the list of allowable currencies for exchange\n", req.ToCurrency)
		return exception.ErrToCurrencyNotAllowed
	}

	// check if toCurrency is same as fromCurrency
	if fromCurrency == req.ToCurrency {
		log.Printf("fromCurrency %s is equal to toCurrency %s\n", fromCurrency, req.ToCurrency)
		return exception.ErrFromCurrencyEqualToCurrency
	}

	// retrieve converted amount and profit
	profit, convertedAmount := utils.CalculateConversionDetails(req.FromAmount, fromCurrency, req.ToCurrency)

	// TODO: add profit into creator's account balance
	fmt.Println("Adding into creator's account balance", profit)

	// retrieve main balance and check if sufficient funds
	allBalances, err := uc.balanceRepository.GetBalances(ctx, tx, userID)
	if err != nil {
		log.Printf("failed to retrieve all balances for user id %d with error: %v\n", userID, err)
		return err
	}

	// convert slice of user balances into map for faster performance of accessing keys in map
	allBalancesMap := make(map[string]float64)
	for _, b := range allBalances {
		allBalancesMap[b.Currency] = b.Balance
	}

	// check if user has sufficient primary balance to perform currency exchange
	if req.FromAmount > allBalancesMap[fromCurrency] {
		log.Printf("user %d has insufficient balance to perform currency exchange", userID)
		return exception.ErrInsufficientFundsForCurrencyExchange
	}

	finalBalancesMap := make(map[string]float64)
	// add into finalBalancesMap on the fromAmount and toAmount
	finalBalancesMap[fromCurrency] = allBalancesMap[fromCurrency] - req.FromAmount

	// check if user has existing toCurrency balance
	if currentValue, found := allBalancesMap[req.ToCurrency]; found {
		finalBalancesMap[req.ToCurrency] = currentValue + convertedAmount
	} else {
		finalBalancesMap[req.ToCurrency] = convertedAmount
	}

	// update user balances
	if err := uc.balanceRepository.UpdateBalances(ctx, tx, userID, finalBalancesMap); err != nil {
		log.Printf("failed to update balances for user id %d with error: %v\n", userID, err)
		return err
	}

	return nil
}
