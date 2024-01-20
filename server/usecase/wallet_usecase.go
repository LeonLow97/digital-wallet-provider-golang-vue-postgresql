package usecase

import (
	"context"
	"log"

	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/LeonLow97/go-clean-architecture/dto"
	"github.com/LeonLow97/go-clean-architecture/exception"
)

type walletUsecase struct {
	walletRepository domain.WalletRepository
}

func NewWalletUsecase(walletRepository domain.WalletRepository) domain.WalletUsecase {
	return &walletUsecase{
		walletRepository: walletRepository,
	}
}

func (uc *walletUsecase) CreateWallet(ctx context.Context, req dto.CreateWalletRequest) error {
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
		Currency: req.Currency,
		TypeID:   walletTypeID,
		UserID:   req.UserID,
	}

	if err := uc.walletRepository.CreateWallet(ctx, &createWallet); err != nil {
		log.Println("error creating wallet", err)
		return err
	}
	return nil
}
