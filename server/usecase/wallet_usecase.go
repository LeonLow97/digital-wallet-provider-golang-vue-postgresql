package usecase

import (
	"context"
	"log"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
	"github.com/jmoiron/sqlx"
)

type walletUsecase struct {
	dbConn            *sqlx.DB
	walletRepository  domain.WalletRepository
	balanceRepository domain.BalanceRepository
}

func NewWalletUsecase(dbConn *sqlx.DB, walletRepository domain.WalletRepository, balanceRepository domain.BalanceRepository) domain.WalletUsecase {
	return &walletUsecase{
		dbConn:            dbConn,
		walletRepository:  walletRepository,
		balanceRepository: balanceRepository,
	}
}

func (uc *walletUsecase) GetWallet(ctx context.Context, userID, walletID int) (*domain.Wallet, error) {
	// retrieve one wallet by user id and wallet ID
	wallet, err := uc.walletRepository.GetWalletByWalletID(ctx, userID, walletID)
	if err != nil {
		log.Printf("failed to get wallet for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	// retrieve wallet balances by user id and wallet id
	walletBalances, err := uc.walletRepository.GetWalletBalancesByUserIDAndWalletID(ctx, userID, walletID)
	if err != nil {
		log.Printf("failed to get wallet balances for user id %d and wallet id %d with error: %v\n", userID, walletID, err)
		return nil, err
	}
	wallet.CurrencyAmount = walletBalances

	return wallet, nil
}

func (uc *walletUsecase) GetWallets(ctx context.Context, userID int) (*[]domain.Wallet, error) {
	// retrieve wallets by user id
	wallets, err := uc.walletRepository.GetWallets(ctx, userID)
	if err != nil {
		log.Printf("failed to get wallets for user id %d with error: %v\n", userID, err)
		return nil, err
	}
	if len(wallets) == 0 {
		return nil, exception.ErrNoWalletsFound
	}

	// retrieve wallet balances by user id
	walletBalances, err := uc.walletRepository.GetWalletBalancesByUserID(ctx, userID)
	if err != nil {
		log.Printf("failed to get wallet balances for user id %d with error: %v\n", userID, err)
		return nil, err
	}

	// walletBalancesMap --> key: walletID, value: { currency: amount }
	walletBalancesMap := make(map[int][]domain.WalletCurrencyAmount)
	for _, wb := range walletBalances {
		if _, found := walletBalancesMap[wb.WalletID]; !found {
			walletBalancesMap[wb.WalletID] = []domain.WalletCurrencyAmount{}
		}
		walletBalancesMap[wb.WalletID] = append(walletBalancesMap[wb.WalletID], wb)
	}

	for idx, w := range wallets {
		if wb, found := walletBalancesMap[w.ID]; found {
			wallets[idx].CurrencyAmount = wb
		}
	}

	return &wallets, nil
}

func (uc *walletUsecase) GetWalletTypes(ctx context.Context) (*[]dto.GetWalletTypesResponse, error) {
	walletTypes, err := uc.walletRepository.GetWalletTypes(ctx)
	if err != nil {
		log.Println("failed to retrieve wallet types with error:", err)
		return nil, err
	}

	return walletTypes, nil
}

func (uc *walletUsecase) CreateWallet(ctx context.Context, userID int, req dto.CreateWalletRequest) error {
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

	// check if user has already created these wallets
	walletValidation, err := uc.walletRepository.PerformWalletValidationByUserID(ctx, userID, req.WalletTypeID)
	if err != nil {
		log.Printf("failed to check wallet exists by user id %d with error: %v\n", userID, err)
		return err
	}
	if walletValidation.WalletExists {
		return exception.ErrWalletAlreadyExists
	}
	if !walletValidation.IsValidWalletType {
		return exception.ErrWalletTypeInvalid
	}

	// retrieve main balance and check if sufficient funds
	allBalances, err := uc.walletRepository.GetAllBalancesByUserID(ctx, userID)
	if err != nil {
		log.Printf("failed to retrieve all balances for user id %d with error: %v\n", userID, err)
		return err
	}

	// convert slice of user balances into map for faster performance of accessing keys in map
	allBalancesMap := make(map[string]float64)
	for _, b := range allBalances {
		allBalancesMap[b.Currency] = b.Balance
	}

	finalBalancesMap := make(map[string]float64)
	currencyAmount := make([]domain.WalletCurrencyAmount, 0)

	// ensure all balances are sufficient to top up new wallet
	for _, a := range req.CurrencyAmount {
		if currentBalance, found := allBalancesMap[a.Currency]; !found {
			// user does not have a balance in this currency
			log.Printf("user %d does not have a balance in this currency\n", userID)
			return exception.ErrBalanceNotFound
		} else {
			if currentBalance < a.Amount {
				log.Printf("user %d has insufficient funds to top up wallet\n", userID)
				return exception.ErrInsufficientFunds
			}

			currencyAmount = append(currencyAmount, domain.WalletCurrencyAmount{
				Amount:   a.Amount,
				Currency: a.Currency,
			})

			finalBalance := allBalancesMap[a.Currency] - a.Amount
			finalBalancesMap[a.Currency] = finalBalance
		}
	}

	// update user balances
	if err := uc.balanceRepository.UpdateBalances(ctx, tx, userID, finalBalancesMap); err != nil {
		log.Printf("failed to update balances for user id %d with error: %v\n", userID, err)
		return err
	}

	// create wallets
	newWallet := &domain.Wallet{
		WalletTypeID: req.WalletTypeID,
		UserID:       userID,
	}
	walletID, err := uc.walletRepository.CreateWallet(ctx, tx, newWallet)
	if err != nil {
		log.Printf("failed to create wallet for user id %d with error: %v\n", userID, err)
		return err
	}

	// insert amount and currency for the wallet
	if err := uc.walletRepository.InsertWalletCurrencyAmount(ctx, tx, walletID, userID, currencyAmount); err != nil {
		log.Printf("failed to insert wallet currency amount for user id %d with error: %v\n", userID, err)
		return err
	}

	return nil
}

func (uc *walletUsecase) UpdateWallet(ctx context.Context, req dto.UpdateWalletRequest) (*dto.UpdateWalletResponse, error) {
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

	// wallet, err := uc.walletRepository.GetWalletByWalletType(ctx, req.UserID, req.Type)
	// if err != nil {
	// 	log.Println("error getting one wallet", err)
	// 	return nil, err
	// }

	// accountBalance, err := uc.walletRepository.GetBalanceByUserID(ctx, req.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// // Check if amount to be updated is negative or positive
	// // negative = withdraw, positive = deposit
	// if req.Balance < 0 {
	// 	if math.Abs(req.Balance) > wallet.Balance {
	// 		return nil, exception.ErrInsufficientFundsForWithdrawal
	// 	}
	// } else if req.Balance > 0 {
	// 	if req.Balance > accountBalance.Balance {
	// 		return nil, exception.ErrInsufficientFundsForDeposit
	// 	}
	// } else {
	// 	// req.Balance is 0, no need to update
	// 	return nil, nil
	// }

	// // update the wallet balance
	// updatedWalletBalance := wallet.Balance + req.Balance
	// wallet.Balance = updatedWalletBalance

	// // TODO: Add concurrency here --> DB Transaction
	// err = uc.walletRepository.UpdateWallet(ctx, wallet)
	// if err != nil {
	// 	return nil, err
	// }

	// updatedAccountBalance := accountBalance.Balance - req.Balance
	// b := domain.Balance{
	// 	Balance:  updatedAccountBalance,
	// 	UserID:   req.UserID,
	// 	Currency: req.BalanceCurrency,
	// }
	// if err := uc.balanceRepository.UpdateBalance(ctx, tx, &b); err != nil {
	// 	return nil, err
	// }

	// return &dto.UpdateWalletResponse{
	// 	WalletID: wallet.ID,
	// 	Type:     wallet.Type,
	// 	Balance:  wallet.Balance,
	// }, nil

	return nil, nil
}
