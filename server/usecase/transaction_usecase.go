package usecase

import (
	"context"
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
	userRepository        domain.UserRepository
	balanceRepository     domain.BalanceRepository
}

func NewTransactionUsecase(dbConn *sqlx.DB, transactionRepo domain.TransactionRepository, walletRepo domain.WalletRepository, userRepo domain.UserRepository, balanceRepo domain.BalanceRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		dbConn:                dbConn,
		transactionRepository: transactionRepo,
		walletRepository:      walletRepo,
		userRepository:        userRepo,
		balanceRepository:     balanceRepo,
	}
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest, userID int) error {
	// ensure amount specified is positive
	if req.Amount <= 0 {
		return exception.ErrAmountMustBePositive
	}

	// Start SQL Transaction, need to lock balance in case use POSTMAN and frontend to update balance at the same time
	tx, err := uc.dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to begin sql transaction in deposit usecase with error: %v\n", err)
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
	if err := uc.transactionRepository.CheckLinkageOfSenderAndBeneficiaryByMobileNumber(ctx, userID, req.MobileNumber); err != nil {
		// if not linked, error will be thrown here
		log.Println("failed to check linkage of sender and beneficiary", err)
		return err
	}

	// retrieve sender details by user id
	userWallet, err := uc.walletRepository.GetUserAndWalletByUserID(ctx, userID, req.SenderWalletID, req.SourceCurrency)
	if err != nil {
		log.Println("failed to get sender wallet details", err)
		return err
	}

	// retrieve beneficiary details by mobile number
	beneficiary, err := uc.userRepository.GetUserAndBalanceByMobileNumber(ctx, req.MobileNumber)
	if err != nil {
		log.Println("failed to get beneficiary details", err)
		return err
	}

	// check if both sender and beneficiary are active
	if !userWallet.Active {
		return exception.ErrUserIsInactive
	}
	if !beneficiary.Active {
		return exception.ErrBeneficiaryIsInactive
	}

	// check if sender id is equal to beneficiary id, cannot send money to yourself
	if userWallet.UserID == beneficiary.ID {
		return exception.ErrUserIDEqualBeneficiaryID
	}

	// check if sender wallet balance is sufficient for transfer and calculate final balance
	if userWallet.Balance < req.Amount {
		return exception.ErrInsufficientFundsInWallet
	}

	// update the balance of both sender and beneficiary
	senderFinalBalance := userWallet.Balance - req.Amount
	beneficiaryFinalBalance := beneficiary.Balance + req.Amount // TODO: perform currency exchange

	updatedSenderWallet := &domain.Wallet{
		Balance:  senderFinalBalance,
		UserID:   userID,
		TypeID:   userWallet.TypeID,
		Currency: userWallet.Currency,
	}

	updatedBeneficiaryBalance := &domain.Balance{
		Balance:  beneficiaryFinalBalance,
		UserID:   beneficiary.ID,
		Currency: req.DestinationCurrency,
	}

	err = uc.walletRepository.UpdateWallet(ctx, updatedSenderWallet)
	if err != nil {
		log.Println("failed to update sender wallet balance", err)
		return err
	}

	err = uc.balanceRepository.UpdateBalance(ctx, tx, updatedBeneficiaryBalance)
	if err != nil {
		log.Println("failed to update beneficiary balance", err)
		return err
	}

	transactionEntity := &domain.Transaction{
		SentAmount:       req.Amount,
		SourceCurrency:   req.SourceCurrency,
		ReceivedAmount:   req.Amount, // TODO: perform currency exchange
		ReceivedCurrency: req.DestinationCurrency,
		SourceOfTransfer: userWallet.Type,
		Status:           utils.SUBMITTED,
	}

	// create 2 transactions in the transactions table for sender and beneficiary
	err = uc.transactionRepository.InsertTransaction(ctx, userID, userID, beneficiary.ID, *transactionEntity)
	if err != nil {
		log.Println("failed to insert transaction", err)
		return err
	}

	return nil
}

func (uc *transactionUsecase) GetTransactions(ctx context.Context, userID int) (*[]domain.Transaction, error) {
	transactions, err := uc.transactionRepository.GetTransactions(ctx, userID)
	if err != nil {
		log.Println("failed to get transactions", err)
		return nil, err
	}

	return transactions, nil
}
