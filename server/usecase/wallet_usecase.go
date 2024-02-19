package usecase

import (
	"context"
	"log"
	"math"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
)

type walletUsecase struct {
	walletRepository  domain.WalletRepository
	balanceRepository domain.BalanceRepository
}

func NewWalletUsecase(walletRepository domain.WalletRepository, balanceRepository domain.BalanceRepository) domain.WalletUsecase {
	return &walletUsecase{
		walletRepository:  walletRepository,
		balanceRepository: balanceRepository,
	}
}

func (uc *walletUsecase) GetWallet(ctx context.Context, userID, walletID int) (*dto.GetWalletResponse, error) {
	wallet, err := uc.walletRepository.GetWalletByWalletID(ctx, userID, walletID)
	if err != nil {
		log.Println("error getting one wallet", err)
		return nil, err
	}

	return &dto.GetWalletResponse{
		WalletID:  wallet.ID,
		Type:      wallet.Type,
		TypeID:    wallet.TypeID,
		Balance:   wallet.Balance,
		Currency:  wallet.Currency,
		CreatedAt: wallet.CreatedAt,
	}, nil
}

func (uc *walletUsecase) GetWallets(ctx context.Context, userID int) (*dto.GetWalletsResponse, error) {
	wallets, err := uc.walletRepository.GetWallets(ctx, userID)
	if err != nil {
		log.Println("error getting wallets", err)
		return nil, err
	}
	if len(wallets) == 0 {
		return nil, exception.ErrNoWalletsFound
	}

	// TODO: might remove this unnecessary for loop and just return domain.Wallet
	resp := &dto.GetWalletsResponse{}
	for _, w := range wallets {
		wallet := dto.GetWalletResponse{
			WalletID:  w.ID,
			Type:      w.Type,
			TypeID:    w.TypeID,
			Balance:   w.Balance,
			Currency:  w.Currency,
			CreatedAt: w.CreatedAt,
		}
		resp.Wallets = append(resp.Wallets, wallet)
	}
	return resp, nil
}

func (uc *walletUsecase) CreateWallet(ctx context.Context, req dto.CreateWalletRequest) error {
	// retrieve main balance and check if sufficient funds
	accountBalance, err := uc.walletRepository.GetBalanceByUserID(ctx, req.UserID)
	if err != nil {
		log.Println("error getting balance by user id", err)
		return err
	}
	if accountBalance.Balance < req.Balance {
		log.Printf("Insufficient Funds; Current Balance: %f, Adding: %f", accountBalance.Balance, req.Balance)
		return exception.ErrInsufficientFunds
	}

	// get wallet types and check if specified wallet type exist
	walletTypes, err := uc.walletRepository.GetWalletTypes(ctx)
	if err != nil {
		log.Println("error getting wallet types", err)
		return err
	}

	var walletTypeID int
	if value, found := walletTypes[req.Type]; !found {
		return exception.ErrWalletTypeInvalid
	} else if found {
		walletTypeID = value
	}

	// check if user already has this wallet created
	isExists, err := uc.walletRepository.CheckWalletExistsByUserID(ctx, req.UserID, req.Type)
	if err != nil {
		log.Println("error checking wallet exists by user id", err)
		return err
	}
	if isExists > 0 {
		return exception.ErrWalletAlreadyExists
	}

	createWallet := domain.Wallet{
		Balance:  req.Balance,
		Currency: req.WalletCurrency,
		TypeID:   walletTypeID,
		UserID:   req.UserID,
	}

	updatedAccountBalance := accountBalance.Balance - req.Balance
	b := domain.Balance{
		Balance:  updatedAccountBalance,
		UserID:   req.UserID,
		Currency: req.BalanceCurrency,
	}
	if err := uc.balanceRepository.UpdateBalance(ctx, &b); err != nil {
		return err
	}

	if err := uc.walletRepository.CreateWallet(ctx, &createWallet); err != nil {
		log.Println("error creating wallet", err)
		return err
	}
	return nil
}

func (uc *walletUsecase) UpdateWallet(ctx context.Context, req dto.UpdateWalletRequest) (*dto.UpdateWalletResponse, error) {
	wallet, err := uc.walletRepository.GetWalletByWalletType(ctx, req.UserID, req.Type)
	if err != nil {
		log.Println("error getting one wallet", err)
		return nil, err
	}

	accountBalance, err := uc.walletRepository.GetBalanceByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Check if amount to be updated is negative or positive
	// negative = withdraw, positive = deposit
	if req.Balance < 0 {
		if math.Abs(req.Balance) > wallet.Balance {
			return nil, exception.ErrInsufficientFundsForWithdrawal
		}
	} else if req.Balance > 0 {
		if req.Balance > accountBalance.Balance {
			return nil, exception.ErrInsufficientFundsForDeposit
		}
	} else {
		// req.Balance is 0, no need to update
		return nil, nil
	}

	// update the wallet balance
	updatedWalletBalance := wallet.Balance + req.Balance
	wallet.Balance = updatedWalletBalance

	// TODO: Add concurrency here --> DB Transaction
	err = uc.walletRepository.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	updatedAccountBalance := accountBalance.Balance - req.Balance
	b := domain.Balance{
		Balance:  updatedAccountBalance,
		UserID:   req.UserID,
		Currency: req.BalanceCurrency,
	}
	if err := uc.balanceRepository.UpdateBalance(ctx, &b); err != nil {
		return nil, err
	}

	return &dto.UpdateWalletResponse{
		WalletID: wallet.ID,
		Type:     wallet.Type,
		Balance:  wallet.Balance,
	}, nil
}
