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

type transactionUsecase struct {
	dbConn                *sqlx.DB
	transactionRepository domain.TransactionRepository
	walletRepository      domain.WalletRepository
	balanceRepository     domain.BalanceRepository
	userRepository        domain.UserRepository
}

func NewTransactionUsecase(dbConn *sqlx.DB, transactionRepo domain.TransactionRepository, walletRepo domain.WalletRepository, balanceRepo domain.BalanceRepository, userRepo domain.UserRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		dbConn:                dbConn,
		transactionRepository: transactionRepo,
		walletRepository:      walletRepo,
		balanceRepository:     balanceRepo,
		userRepository:        userRepo,
	}
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID int) error {
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

	// check if sender id is linked to beneficiary id
	beneficiaryID, isBeneficiaryActive, isMFAConfigured, err := uc.transactionRepository.CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx, userID, req.BeneficiaryMobileCountryCode, req.BeneficiaryMobileNumber)
	if err != nil {
		// if not linked, error will be thrown here
		log.Println("failed to check linkage of sender and beneficiary", err)
		return err
	}
	if !isBeneficiaryActive {
		return exception.ErrBeneficiaryIsInactive
	}
	if !isMFAConfigured {
		return exception.ErrBeneficiaryMFANotConfigured
	}

	// check if sender id is equal to beneficiary id
	if userID == beneficiaryID {
		log.Println("sender id cannot be equal to beneficiary id when performing transaction")
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// check if sender wallet id is linked to user id
	isValidSenderWallet, walletName, err := uc.transactionRepository.CheckValidityOfSenderIDAndWalletID(ctx, userID, req.SenderWalletID)
	if err != nil {
		log.Printf("failed to validate sender wallet id %d with error: %v\n", req.SenderWalletID, err)
		return err
	}
	if !isValidSenderWallet {
		log.Println("sender wallet is invalid, forbid request")
		return exception.ErrSenderWalletInvalid
	}

	// retrieve wallet details and wallet balances for sender
	senderWalletBalances, err := uc.walletRepository.GetWalletBalancesByUserIDAndWalletID(ctx, userID, req.SenderWalletID)
	if err != nil {
		log.Printf("failed to retrieve wallet balances for user id %d and wallet id %d with error %v\n", userID, req.SenderWalletID, err)
		return err
	}

	// TODO: allow internal transfer in same wallet
	// Using a map to allow internal transfer in the wallet if the source currency does not have enough funds
	senderWalletBalancesMap := make(map[string]float64)
	for _, b := range senderWalletBalances {
		senderWalletBalancesMap[b.Currency] = b.Amount
	}

	// TODO: allow internal transfer in same wallet
	// check if sender wallet balance if sufficient for transfer
	if senderWalletBalancesMap[req.SourceCurrency] < req.SourceAmount {
		return exception.ErrInsufficientFundsInWallet
	}

	// retrieve beneficiary balances by beneficiary ID (equivalent to userID for beneficiary)
	beneficiaryBalances, err := uc.balanceRepository.GetBalances(ctx, tx, beneficiaryID)
	if err != nil {
		log.Printf("failed to retrieve balances for beneficiary id %d with error: %v\n", beneficiaryID, err)
		return err
	}

	beneficiaryBalancesMap := make(map[string]float64)
	for _, b := range beneficiaryBalances {
		beneficiaryBalancesMap[b.Currency] = b.Balance
	}

	// update balance of sender wallet
	finalSenderWalletBalancesMap := make(map[string]float64)
	finalSenderWalletBalancesMap[req.SourceCurrency] = senderWalletBalancesMap[req.SourceCurrency] - req.SourceAmount

	// update balance of beneficiary
	var finalDestinationAmount float64
	var finalDestinationCurrency string

	if _, found := beneficiaryBalancesMap[req.SourceCurrency]; found {
		// beneficiary has balance of the same currency as source currency
		beneficiaryBalancesMap[req.SourceCurrency] += req.SourceAmount

		finalDestinationAmount = req.SourceAmount
		finalDestinationCurrency = req.SourceCurrency
	} else {
		mainDestinationCurrency := utils.CountryCodeToCurrencyMap[req.BeneficiaryMobileCountryCode]
		profit, transferAmount := utils.CalculateConversionDetails(req.SourceAmount, req.SourceCurrency, mainDestinationCurrency)

		// TODO: Capture profit and add to creator's account
		fmt.Println("Profit for exchange -->", profit)

		beneficiaryBalancesMap[mainDestinationCurrency] += transferAmount

		finalDestinationAmount = transferAmount
		finalDestinationCurrency = mainDestinationCurrency
	}

	// update sender wallet balances
	if err := uc.walletRepository.CashOutWalletBalances(ctx, tx, userID, req.SenderWalletID, finalSenderWalletBalancesMap); err != nil {
		log.Printf("failed to cash out wallet balances for sender id %d with error: %v\n", userID, err)
		return err
	}

	// update beneficiary balances
	if err := uc.balanceRepository.UpdateBalances(ctx, tx, beneficiaryID, beneficiaryBalancesMap); err != nil {
		log.Printf("failed to update beneficiary id %d balances with error: %v\n", beneficiaryID, err)
		return err
	}

	// create transaction for sender wallet
	sender, err := uc.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve sender by user id %d with error: %v\n", userID, err)
		return err
	}

	// transaction entity
	transactionEntity := &domain.Transaction{
		SenderID:                userID,
		BeneficiaryID:           beneficiaryID,
		SenderMobileNumber:      fmt.Sprintf("%s %s", sender.MobileCountryCode, sender.MobileNumber),
		BeneficiaryMobileNumber: fmt.Sprintf("%s %s", req.BeneficiaryMobileCountryCode, req.BeneficiaryMobileNumber),
		SourceAmount:            req.SourceAmount,
		SourceCurrency:          req.SourceCurrency,
		DestinationAmount:       finalDestinationAmount,
		DestinationCurrency:     finalDestinationCurrency,
		Status:                  "COMPLETED",
	}

	// create transaction for sender
	transactionEntity.SourceOfTransfer = walletName
	if err := uc.transactionRepository.InsertTransaction(ctx, tx, userID, *transactionEntity); err != nil {
		log.Printf("failed to create transaction for sender ID %d with error: %v\n", userID, err)
		return err
	}

	// create transaction for beneficiary
	transactionEntity.SourceOfTransfer = fmt.Sprintf("Main Balance %s", finalDestinationCurrency)
	if err := uc.transactionRepository.InsertTransaction(ctx, tx, beneficiaryID, *transactionEntity); err != nil {
		log.Printf("failed to create transaction for beneficiary ID %d with error: %v\n", beneficiaryID, err)
		return err
	}

	return nil
}

func (uc *transactionUsecase) GetTransactions(ctx context.Context, userID int) (*[]domain.Transaction, error) {
	transactions, err := uc.transactionRepository.GetTransactions(ctx, userID)
	if err != nil {
		log.Printf("failed to get transactions for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	return transactions, nil
}
